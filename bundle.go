package bono

import (
	"errors"
	"log"
	"net"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

type options struct {
	Adapter string
}

var Options *options = &options{
	Adapter: "fh",
}

type Next func() error

type Middleware func(c *Context, next Next) error

type Bundle struct {
	middlewares []Middleware
}

func (b *Bundle) Use(m Middleware) *Bundle {
	b.middlewares = append(b.middlewares, m)
	return b

}

func (b *Bundle) Listen(address string, reuseaddr bool) error {
	switch Options.Adapter {
	case "fh":
		var (
			listener net.Listener
			err      error
		)
		if reuseaddr {
			if listener, err = reuseport.Listen("tcp4", address); err != nil {
				return err
			}

			return fasthttp.Serve(listener, b.FasthttpCallback)
		} else {
			return fasthttp.ListenAndServe(address, b.FasthttpCallback)
		}
	default:
		return errors.New("Unimplemented yet! " + Options.Adapter)
	}

	return nil
}

func (b *Bundle) Dispatch(context *Context) {
	err := b.dispatchMiddleware(0, context)
	if err != nil {
		switch err.Error() {
		case "Delegated":
			return
		case "Stop":
			log.Printf("Caught error: %s", err.Error())
			context.SetStatus(500)
			context.SetBody([]byte(err.Error() + "\n"))
		}
	}

	if context.Status() == 404 && context.Body() != nil {
		context.SetStatus(200)
	}
}

func (b *Bundle) dispatchMiddleware(i int, context *Context) error {
	if len(b.middlewares) > i {
		middleware := b.middlewares[i]
		return middleware(context, func() error {
			return b.dispatchMiddleware(i+1, context)
		})
	}
	return nil
}

func New() *Bundle {
	return &Bundle{}
}
