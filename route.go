package bono

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type (
	Next func() error

	Middleware func(c *Context, next Next) error

	Route func(c *Context) (interface{}, error)

	RouteSignature struct {
		Uri     string
		Method  string
		Methods []string
		Handler Route

		kind     byte
		pattern  *regexp.Regexp
		args     []string
		uriBytes []byte
	}

	Router struct {
		Bundle Bundle

		HEAD   []*RouteSignature
		GET    []*RouteSignature
		PUT    []*RouteSignature
		POST   []*RouteSignature
		DELETE []*RouteSignature
	}
)

/**
 * Router methods
 */

func (r *Router) findRouteSignature(method []byte, path []byte) *RouteSignature {
	var routes []*RouteSignature

	switch string(method) {
	case "HEAD":
		routes = r.HEAD
	case "GET":
		routes = r.GET
	case "POST":
		routes = r.POST
	case "PUT":
		routes = r.PUT
	case "DELETE":
		routes = r.DELETE
	}

	for _, rs := range routes {
		if rs.Satisfy(path) {
			return rs
		}
	}
	return nil
}

func (r *Router) routeRegexp(uri []byte) (*regexp.Regexp, []string, error) {

	chunks := bytes.Split(uri, []byte{'['})
	if len(chunks) > 2 {
		return nil, nil, errors.New("Invalid use of optional params")
	}

	extractorRe := regexp.MustCompile("^{([^}]+)}$")
	replacerRe := regexp.MustCompile("{([^}]+)}")

	var tokens []string

	result := replacerRe.ReplaceAllStringFunc(string(chunks[0]), func(token string) string {
		result := extractorRe.FindStringSubmatch(token)
		tokens = append(tokens, result[1])
		return `([^\/]+)`
	})

	// TODO optional param is not supported yet

	re, err := regexp.Compile(result)
	return re, tokens, err
}

func (r *Router) Map(signature *RouteSignature) error {
	if signature.Uri == "" {
		return errors.New("Uri is undefined")
	}
	if signature.Methods == nil {
		if signature.Method == "" {
			signature.Methods = []string{"GET"}
		} else {
			signature.Methods = []string{signature.Method}
		}
	}
	signature.Method = ""

	if !signature.IsStatic() {
		re, tokens, err := r.routeRegexp(signature.UriBytes())
		if err != nil {
			return err
		}

		signature.pattern = re
		signature.args = tokens
	}

	for _, method := range signature.Methods {
		s := *signature
		signatureClone := &s
		signatureClone.Method = strings.ToUpper(method)

		switch signatureClone.Method {
		case "HEAD":
			r.HEAD = append(r.HEAD, signatureClone)
		case "GET":
			r.GET = append(r.GET, signatureClone)
		case "POST":
			r.POST = append(r.POST, signatureClone)
		case "PUT":
			r.PUT = append(r.PUT, signatureClone)
		case "DELETE":
			r.DELETE = append(r.DELETE, signatureClone)
		}
	}

	return nil
}

func (r *Router) Route(ctx *Context, next Next) error {
	//satisfy route
	rs := r.findRouteSignature(ctx.Method(), ctx.Path())
	if rs != nil {
		ctx.SetStatus(200)

		rs.fetchAttributes(ctx)
		ctx.Attr()["route.uri"] = rs.Uri

		state, err := rs.Handler(ctx)
		if err != nil {
			return err
		}
		if state != nil {
			ctx.SetState(state)
		}

		return nil
	}

	return next()

}

/**
 * RouteSignature methods
 */

func (rs *RouteSignature) IsStatic() bool {
	if rs.kind == 0 {
		hasRegex, _ := regexp.MatchString("[[{]", rs.Uri)
		if hasRegex {
			rs.kind = 'v'
		} else {
			rs.kind = 's'
		}
	}
	return rs.kind == 's'
}

func (rs *RouteSignature) Satisfy(path []byte) bool {
	if rs.IsStatic() {
		return bytes.Equal(rs.UriBytes(), path)
	} else {
		return rs.pattern.Match(path)
	}
}

func (rs *RouteSignature) UriBytes() []byte {
	if rs.uriBytes == nil {
		rs.uriBytes = []byte(rs.Uri)
	}

	return rs.uriBytes
}

func (rs *RouteSignature) fetchAttributes(ctx *Context) {
	if rs.kind != 'v' {
		return
	}
	matches := rs.pattern.FindAllSubmatch(ctx.Path(), -1)
	for i, arg := range rs.args {
		ctx.Attr()[arg] = matches[0][i+1]
	}
}
