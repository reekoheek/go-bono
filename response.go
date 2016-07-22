package bono

import "strings"

type (
	Response interface {
		Status() int
		SetStatus(status int) error
		Body() []byte
		SetBody(body []byte) error
		Headers() map[string]string
		Set(key string, value string)
		SetContentType(contentType string)
	}

	ResponseImpl struct {
		status  int
		body    []byte
		headers map[string]string
	}
)

func (r *ResponseImpl) Status() int {
	if r.status == 0 {
		r.status = 404
	}
	return r.status
}

func (r *ResponseImpl) SetStatus(status int) error {
	r.status = status
	return nil
}

func (r *ResponseImpl) Body() []byte {
	return r.body
}

func (r *ResponseImpl) SetBody(body []byte) error {
	r.body = body
	return nil
}

func (r *ResponseImpl) Headers() map[string]string {
	return r.headers
}

func (r *ResponseImpl) Set(key string, value string) {
	if r.headers == nil {
		r.headers = map[string]string{}
	}
	r.headers[strings.ToLower(key)] = value
}

func (r *ResponseImpl) SetContentType(contentType string) {
	r.Set("content-type", contentType)
}
