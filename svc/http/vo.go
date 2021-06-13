package ddhttp

import "net/http"

//go:generate go-doudou name --file $GOFILE -o

// POST /usersvc/pageusers HTTP/1.1
//Host: localhost:6060
//Content-Length: 80
//Content-Type: application/json
//User-Agent: go-resty/2.6.0 (https://github.com/go-resty/resty)
//X-Request-Id: d1e4dc83-18be-493e-be5b-2e0faaca90ec
//
//{"filter":{"dept":99,"name":"Jack"},"page":{"orders":null,"pageNo":2,"size":10}}
type HttpLog struct {
	ClientIp          string
	HttpMethod        string
	Uri               string
	Proto             string
	Host              string
	ReqContentLength  int64
	ReqHeader         http.Header
	RequestId         string
	RawReq            string
	RespBody          string
	StatusCode        int
	RespHeader        http.Header
	RespContentLength int
	ElapsedTime       string
	// in ms
	Elapsed int64
}
