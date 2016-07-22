package bono

type (
	Context struct {
		Request  Request
		Response Response
		state    interface{}
	}
)

func (c *Context) State() interface{} {
	return c.state
}

func (c *Context) SetState(state interface{}) *Context {
	c.state = state
	return c
}

func (c *Context) Status() int {
	return c.Response.Status()
}

func (c *Context) SetStatus(status int) error {
	return c.Response.SetStatus(status)
}

func (c *Context) Body() []byte {
	return c.Response.Body()
}

func (c *Context) SetBody(body []byte) error {
	return c.Response.SetBody(body)
}

func (c *Context) Method() []byte {
	return c.Request.Method()
}

func (c *Context) SetMethod(method []byte) error {
	return c.Request.SetMethod(method)
}

func (c *Context) Path() []byte {
	return c.Request.Path()
}

func (c *Context) Base() []byte {
	return c.Request.Base()
}

func (c *Context) Attr() map[string]interface{} {
	return c.Request.Attr()
}

func (c *Context) Redirect(url string, status ...int) error {
	if status == nil {
		status = append(status, 302)
	}
	c.SetStatus(status[0])
	c.Set("Location", url)
	return Stop
}

func (c *Context) Shift(uri []byte) *Context {
	c.Request.Shift(uri)
	return c
}

func (c *Context) Unshift(uri []byte) *Context {
	c.Request.Unshift(uri)
	return c
}

func (c *Context) Set(key string, value string) *Context {
	c.Response.Set(key, value)
	return c
}

func (c *Context) SetContentType(contentType string) *Context {
	c.Response.SetContentType(contentType)
	return c
}

func (c *Context) ParseBody() interface{} {
	return c.Request.ParseBody()
}
