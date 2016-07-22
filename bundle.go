package bono

import (
	"bytes"
	"errors"
	"log"
	"net"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

type (
	Bundle interface {
		Use(m Middleware) Bundle
		Finalize() Bundle
		Listen(address string, reuseaddr bool) error
		Dispatch(context *Context) error
		AddBundle(signature *BundleSignature) Bundle
		RouteMap(signature *RouteSignature) Bundle
		Bundles() []*BundleSignature
		Router() *Router
	}

	BundleImpl struct {
		middlewares []Middleware
		bundles     []*BundleSignature
		router      *Router
	}

	BundleSignature struct {
		Uri      string
		Handler  Bundle
		uriBytes []byte
	}

	BundleOptions struct {
		Middlewares []Middleware
		Routes      []*RouteSignature
		Bundles     []*BundleSignature
	}
)

/**
 * Bundle methods
 */
func (b *BundleImpl) Use(m Middleware) Bundle {
	b.middlewares = append(b.middlewares, m)
	return b
}

func (b *BundleImpl) Finalize() Bundle {
	b.Use(b.internalMiddleware)
	return b
}

func (b *BundleImpl) Listen(address string, reuseaddr bool) error {
	b.Finalize()
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

func (b *BundleImpl) Dispatch(context *Context) error {
	if err := b.dispatchMiddleware(0, context); err != nil {
		switch err {
		case Stop:
		default:
			log.Printf("Caught error: %s", err.Error())
			context.SetStatus(500)
			context.SetState(err.Error())
		}
		return err
	}

	if context.Status() == 404 && context.Body() != nil {
		context.SetStatus(200)
	}

	return nil
}

func (b *BundleImpl) AddBundle(signature *BundleSignature) Bundle {
	signature.Handler.Finalize()
	b.bundles = append(b.bundles, signature)
	return b
}

func (b *BundleImpl) RouteMap(signature *RouteSignature) Bundle {
	if err := b.router.Map(signature); err != nil {
		panic(err.Error())
	}
	return b
}

func (b *BundleImpl) Bundles() []*BundleSignature {
	return b.bundles
}

func (b *BundleImpl) Router() *Router {
	return b.router
}

func (b *BundleImpl) dispatchMiddleware(i int, context *Context) error {
	if len(b.middlewares) > i {
		middleware := b.middlewares[i]
		return middleware(context, func() error {
			return b.dispatchMiddleware(i+1, context)
		})
	}
	return nil
}

func (b *BundleImpl) findBundleSignature(path []byte) *BundleSignature {
	for _, bs := range b.Bundles() {
		if bytes.HasPrefix(path, bs.UriBytes()) {
			return bs
		}
	}
	return nil
}

func (b *BundleImpl) internalMiddleware(ctx *Context, next Next) error {
	// satisfy bundle
	bs := b.findBundleSignature(ctx.Path())
	if bs != nil {
		ctx.Attr()["route.bundle"] = bs.Handler
		ctx.Shift(bs.UriBytes())
		err := bs.Handler.Dispatch(ctx)
		ctx.Unshift(bs.UriBytes())
		return err
	}

	return b.router.Route(ctx, next)
}

/**
 * BundleSignature methods
 */

func (bs *BundleSignature) UriBytes() []byte {
	if bs.uriBytes == nil {
		bs.uriBytes = []byte(bs.Uri)
	}
	return bs.uriBytes
}

func New(options *BundleOptions) Bundle {
	bundle := &BundleImpl{
		router: &Router{},
	}

	if options != nil {
		for _, m := range options.Middlewares {
			bundle.Use(m)
		}

		for _, b := range options.Bundles {
			bundle.AddBundle(b)
		}

		for _, r := range options.Routes {
			bundle.RouteMap(r)
		}
	}

	return bundle
}
