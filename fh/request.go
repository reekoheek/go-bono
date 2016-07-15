package fh

import "github.com/valyala/fasthttp"

type Request struct {
	Context *fasthttp.RequestCtx
	method  string
	path    string
}

func (r *Request) Method() string {
	if r.method == "" {
		r.method = string(r.Context.Method())
	}
	return r.method
}

// func (r *Request) SetMethod(method string) error {
// 	r.method = method
// 	return nil
// }

func (r *Request) Path() string {
	if r.path == "" {
		r.path = string(r.Context.Path())
	}
	return r.path
}

// func (r *Request) SetPath(path string) error {
// 	r.path = path
// 	return nil
// }
