package bono

type Context struct {
	Request  Request
	Response Response
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

func (c *Context) Method() string {
	return c.Request.Method()
}

// func (c *Context) SetMethod(method string) error {
// 	return c.Request.SetMethod(method)
// }

func (c *Context) Path() string {
	return c.Request.Path()
}

// func (c *Context) SetPath(path string) error {
// 	return c.Request.SetPath(path)
// }

// type Context struct {
// 	Request    *http.Request
// 	Response   *Response
// 	Attributes map[string]interface{}
// }
//
// func (c *Context) Set(key string, value interface{}) {
// 	if c.Attributes == nil {
// 		c.Attributes = make(map[string]interface{})
// 	}
// 	c.Attributes[key] = value
// }
//
// func (c *Context) Get(key string) interface{} {
// 	return c.Attributes[key]
// }
//
// func (c *Context) Redirect(url string, status ...int) error {
// 	if len(status) == 0 {
// 		status = append(status, 302)
// 	}
// 	c.Response.Status = status[0]
// 	//c.Response.Writer.Header().Set("Location", url)
// 	return errors.New("Stop")
// }
