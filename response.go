package curl

import (
	"io/ioutil"
	"net/http"
)

type Response struct {
	Raw     *http.Response    // http response 內容
	Headers map[string]string // 整理後的 header 內容
	Body    []byte            // http response body 內容
}

func NewResponse() *Response {
	return &Response{}
}

// isOk check http response statuscode
func (r *Response) isOk() bool {
	// 防呆用
	if r == nil || r.Raw == nil {
		return false
	}

	return r.Raw.StatusCode < 300
}

// parseHeaders Handle http response header
func (r *Response) parseHeaders() {
	r.Headers = make(map[string]string)
	for k, v := range r.Raw.Header {
		r.Headers[k] = v[0]
	}
}

// parseBody Handle http response body
func (r *Response) parseBody() (err error) {
	if r.Body, err = ioutil.ReadAll(r.Raw.Body); err != nil {
		return
	}

	return
}
