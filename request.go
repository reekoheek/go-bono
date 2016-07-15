package bono

type Request interface {
	Method() string
	// SetMethod(method string) error
	Path() string
	// SetPath(path string) error
}
