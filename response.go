package bono

type Response interface {
	Status() int
	SetStatus(status int) error
	Body() []byte
	SetBody(body []byte) error
}

type ResponseImpl struct {
	status int
	body   []byte
}

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
