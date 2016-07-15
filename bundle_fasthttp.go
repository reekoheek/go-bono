package bono

import (
	"github.com/reekoheek/go-bono/fh"
	"github.com/valyala/fasthttp"
)

func (b *Bundle) FasthttpCallback(backedCtx *fasthttp.RequestCtx) {
	ctx := &Context{
		Request: &fh.Request{
			Context: backedCtx,
		},
		Response: &ResponseImpl{},
	}
	b.Dispatch(ctx)
	backedCtx.SetStatusCode(ctx.Status())
	backedCtx.SetBody(ctx.Body())
}
