package ddhttp

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/unionj-cloud/go-doudou/framework/cache"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/framework/registry/nacos"
	"github.com/wubin1989/nacos-sdk-go/vo"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Headers borrowed from labstack/echo
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderReferrerPolicy                  = "Referrer-Policy"
)

type ProxyTarget struct {
	Name string
	URL  *url.URL
}

type ProxyConfig struct {
	ProviderStore cache.IStore
	// To customize the transport to remote.
	// Examples: If custom TLS certificates are required.
	Transport http.RoundTripper

	// ModifyResponse defines function to modify response from ProxyTarget.
	ModifyResponse func(*http.Response) error
}

func captureTokens(pattern *regexp.Regexp, input string) *strings.Replacer {
	groups := pattern.FindAllStringSubmatch(input, -1)
	if groups == nil {
		return nil
	}
	values := groups[0][1:]
	replace := make([]string, 2*len(values))
	for i, v := range values {
		j := 2 * i
		replace[j] = "$" + strconv.Itoa(i+1)
		replace[j+1] = v
	}
	return strings.NewReplacer(replace...)
}

func getPath(r *http.Request) string {
	path := r.URL.RawPath
	if path == "" {
		path = r.URL.Path
	}
	return path
}

func isWebSocket(r *http.Request) bool {
	upgrade := r.Header.Get(HeaderUpgrade)
	return strings.ToLower(upgrade) == "websocket"
}

func Proxy(proxyConfig ProxyConfig) func(inner http.Handler) http.Handler {
	if proxyConfig.ProviderStore == nil {
		arc, _ := lru.NewARC(128)
		proxyConfig.ProviderStore = arc
	}
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isWebSocket(r) || r.Header.Get(HeaderAccept) == "text/event-stream" {
				http.Error(w, fmt.Sprintf("not support"), http.StatusBadGateway)
				return
			}
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) <= 1 {
				http.Error(w, fmt.Sprintf("request url must be prefixed / + service name"), http.StatusBadGateway)
				return
			}
			serviceName := parts[1]
			modes := strings.Split(os.Getenv("GDD_SERVICE_DISCOVERY_MODE"), ",")
			var provider registry.IServiceProvider
			for _, mode := range modes {
				switch mode {
				case "nacos":
					cluster := config.GddNacosClusterName.LoadOrDefault(config.DefaultGddNacosClusterName)
					group := config.GddNacosGroupName.LoadOrDefault(config.DefaultGddNacosGroupName)
					_, err := nacos.NamingClient.GetService(vo.GetServiceParam{
						Clusters:    []string{cluster},
						ServiceName: serviceName,
						GroupName:   group,
					})
					if err != nil {
						continue
					}
					if value, ok := proxyConfig.ProviderStore.Get(serviceName); ok {
						if provider, ok = value.(*NacosWRRServiceProvider); ok {
							break
						}
					}
					provider = NewNacosWRRServiceProvider(serviceName, WithNacosClusters([]string{cluster}), WithNacosGroupName(group))
					proxyConfig.ProviderStore.Add(serviceName, provider)
				default:
					nodes, err := registry.AllNodes()
					if err != nil {
						continue
					}
					exist := false
					for _, node := range nodes {
						if registry.SvcName(node) == serviceName {
							exist = true
							break
						}
					}
					if !exist {
						continue
					}
					if value, ok := proxyConfig.ProviderStore.Get(serviceName); ok {
						if provider, ok = value.(*SmoothWeightedRoundRobinProvider); ok {
							break
						}
					}
					provider = NewSmoothWeightedRoundRobinProvider(serviceName)
					proxyConfig.ProviderStore.Add(serviceName, provider)
				}
				if provider != nil {
					break
				}
			}
			if provider == nil {
				http.Error(w, fmt.Sprintf("available server for service %s not found", serviceName), http.StatusBadGateway)
				return
			}
			k := regexp.MustCompile(strings.Replace(fmt.Sprintf("/%s/*", serviceName), "*", "(\\S*)", -1))
			replacer := captureTokens(k, getPath(r))
			if replacer != nil {
				r.URL.Path = replacer.Replace("/$1")
			}
			parsed, err := url.Parse(provider.SelectServer())
			if err != nil {
				http.Error(w, fmt.Sprintf("available server for service %s not found with error: %s", serviceName, err), http.StatusBadGateway)
				return
			}
			tgt := &ProxyTarget{
				Name: serviceName,
				URL:  parsed,
			}
			proxyHTTP(tgt, w, proxyConfig).ServeHTTP(w, r)
		})
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func proxyHTTP(tgt *ProxyTarget, w http.ResponseWriter, config ProxyConfig) http.Handler {
	target := tgt.URL
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		req.Header.Set("Host", target.Host)
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		desc := target.String()
		if tgt.Name != "" {
			desc = fmt.Sprintf("%s(%s)", tgt.Name, tgt.URL.String())
		}
		http.Error(w, fmt.Sprintf("remote %s unreachable, could not forward: %v", desc, err), http.StatusBadGateway)
	}
	proxy.Transport = config.Transport
	proxy.ModifyResponse = config.ModifyResponse
	return proxy
}
