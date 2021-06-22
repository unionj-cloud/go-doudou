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
	ClientIp          string      `json:"clientIp,omitempty"`
	HttpMethod        string      `json:"httpMethod,omitempty"`
	Uri               string      `json:"uri,omitempty"`
	Proto             string      `json:"proto,omitempty"`
	Host              string      `json:"host,omitempty"`
	ReqContentLength  int64       `json:"reqContentLength,omitempty"`
	ReqHeader         http.Header `json:"reqHeader,omitempty"`
	RequestId         string      `json:"requestId,omitempty"`
	RawReq            string      `json:"rawReq,omitempty"`
	RespBody          string      `json:"respBody,omitempty"`
	StatusCode        int         `json:"statusCode,omitempty"`
	RespHeader        http.Header `json:"respHeader,omitempty"`
	RespContentLength int         `json:"respContentLength,omitempty"`
	ElapsedTime       string      `json:"elapsedTime,omitempty"`
	// in ms
	Elapsed int64 `json:"elapsed,omitempty"`
}
