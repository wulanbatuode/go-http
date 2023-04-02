package handmadehttp

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isSameFunc(t *testing.T, expect, actual HandlerFunc) bool {
	t.Helper()
	return reflect.ValueOf(expect) == reflect.ValueOf(actual)
}

func assertEqualFunc(t *testing.T, isSame bool, expectFn, actualFn HandlerFunc) {
	t.Helper()
	assert.Equal(t, isSame, isSameFunc(t, expectFn, actualFn))
}

func TestMultiplexer(t *testing.T) {
	m := NewMultiplexer(nil)
	assertEqualFunc(t, true, NotFoundHandler, *m.findHandler("/"))
	testCases := []struct {
		URI string
		Fn  HandlerFunc
	}{
		{"/", func(req *Request, res *Response) error { return nil }},
		{"/a", func(req *Request, res *Response) error { return nil }},
		{"/b", func(req *Request, res *Response) error { return nil }},
		{"/c", func(req *Request, res *Response) error { return nil }},
		{"/d/e/f", func(req *Request, res *Response) error { return nil }},
	}
	for _, v := range testCases {
		m.UpdateHandler(v.URI, v.Fn)
	}
	for _, v := range testCases {
		assertEqualFunc(t, true, v.Fn, *m.findHandler(v.URI))
		assertEqualFunc(t, false, NotFoundHandler, *m.findHandler(v.URI))
	}
	assertEqualFunc(t, false, NotFoundHandler, *m.findHandler("/"))

}
