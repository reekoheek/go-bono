package bono

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

type (
	Request interface {
		Method() []byte
		SetMethod(method []byte) error
		Path() []byte
		Base() []byte
		Shift(uri []byte)
		Unshift(uri []byte)
		Attr() map[string]interface{}
		ParseBody() *fasthttp.Args
		FetchBody(interface{}) error
		QueryParams() *fasthttp.Args
	}

	RequestImpl struct {
		Context *fasthttp.RequestCtx
		method  []byte
		path    []byte
		base    []byte
		attr    map[string]interface{}
		headers map[string]string
	}
)

func (r *RequestImpl) Method() []byte {
	if r.method == nil {
		r.method = r.Context.Method()
	}
	return r.method
}

func (r *RequestImpl) SetMethod(method []byte) error {
	r.method = method
	return nil
}

func (r *RequestImpl) Path() []byte {
	if r.path == nil {
		r.path = r.Context.Path()
	}
	return r.path
}

func (r *RequestImpl) SetPath(path []byte) {
	r.path = path
}

func (r *RequestImpl) Base() []byte {
	if r.base == nil {
		r.base = []byte{'/'}
	}
	return r.base
}

func (r *RequestImpl) Shift(uri []byte) {
	uriLen := len(uri)
	r.path = r.path[uriLen:]
	if len(r.path) == 0 {
		r.path = []byte("/")
	}
	if len(r.base) == 1 {
		r.base = uri
	} else {
		r.base = append(r.base, uri...)
	}
}

func (r *RequestImpl) Unshift(uri []byte) {
	baseLen := len(r.base)
	uriLen := len(uri)
	r.path = append(uri, r.path...)
	if baseLen > 1 {
		r.base = r.base[:baseLen-uriLen]
	}
}

func (r *RequestImpl) Attr() map[string]interface{} {
	if r.attr == nil {
		r.attr = map[string]interface{}{}
	}
	return r.attr
}

func (r *RequestImpl) ParseBody() *fasthttp.Args {
	return r.Context.PostArgs()
}

func (r *RequestImpl) FetchBody(model interface{}) error {
	args := r.Context.PostArgs()
	val := reflect.ValueOf(model).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		field := tag.Get("field")
		if field == "" {
			field = strings.ToLower(typeField.Name)
		}

		// fmt.Printf("> kind: %10s %t field: %10s value: %10s\n", val.Field(i).Kind(), val.Field(i).CanSet(), field, args.Peek(field))

		strValue := string(args.Peek(field))
		f := val.Field(i)
		switch f.Kind() {
		case reflect.Bool:
			v, _ := strconv.ParseBool(strValue)
			f.SetBool(v)
		case reflect.Int:
			v, _ := strconv.ParseInt(strValue, 10, 64)
			f.SetInt(v)
		case reflect.Float32, reflect.Float64:
			v, _ := strconv.ParseFloat(strValue, 64)
			f.SetFloat(v)
		case reflect.String:
			f.SetString(strValue)
		}
	}
	return nil
}

func (r *RequestImpl) QueryParams() *fasthttp.Args {
	return r.Context.QueryArgs()
}
