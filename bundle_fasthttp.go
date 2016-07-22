package bono

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func (b *BundleImpl) FasthttpCallback(backedCtx *fasthttp.RequestCtx) {
	ctx := &Context{
		Request: &RequestImpl{
			Context: backedCtx,
		},
		Response: &ResponseImpl{},
	}
	switch err := b.Dispatch(ctx); err {
	case Halt:
		return
	default:
		if ctx.Status() == 404 && len(ctx.Body()) == 0 {
			ctx.SetBody([]byte(fmt.Sprintf("%d", ctx.Status())))
		}

		backedCtx.SetStatusCode(ctx.Status())
		backedCtx.SetBody(ctx.Body())

		// log.Println("xxxx", ctx.Response.Headers())

		for k, v := range ctx.Response.Headers() {
			backedCtx.Response.Header.Set(k, v)
		}
	}
}
