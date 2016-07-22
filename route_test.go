package bono

import (
	"log"
	"regexp"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestRouter_routeRegexp_withSingleVar(t *testing.T) {
	var (
		router *Router = &Router{}
		uri    []byte  = []byte("/{id}")
		re     *regexp.Regexp
		tokens [][]byte
		err    error
	)

	if re, tokens, err = router.routeRegexp(uri); err != nil {
		t.Error(err.Error())
		return
	}

	if re == nil {
		t.Error("routeRegexp must return Regexp struct")
		return
	}

	if len(tokens) < 1 {
		t.Error("routeRegexp must return tokens")
		return
	}
}

func TestRouter_routeRegexp_withMultipleVars(t *testing.T) {
	var (
		router *Router = &Router{}
		uri    []byte  = []byte("/foo/{fooId}/bar/{id}")
		re     *regexp.Regexp
		tokens [][]byte
		err    error
	)

	if re, tokens, err = router.routeRegexp(uri); err != nil {
		t.Error(err.Error())
		return
	}

	if re == nil {
		t.Error("routeRegexp must return Regexp struct")
		return
	}

	if len(tokens) < 1 {
		t.Error("routeRegexp must return tokens")
		return
	}
}
