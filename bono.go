package bono

import "errors"

type (
	typeOptions struct {
		Adapter string
	}
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
)

var (
	Options *typeOptions = &typeOptions{
		Adapter: "fh",
	}
	// halt will immediate return without set response data
	Halt error = errors.New("!!halt")
	// stop will return then set response data
	Stop error = errors.New("!!stop")
)
