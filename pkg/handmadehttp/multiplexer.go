package handmadehttp

import "strings"

type HandlerFunc func(req *Request, res *Response) error

func NotFoundHandler(req *Request, res *Response) error {
	res.StatusCode = 404
	res.SetContent([]byte{})
	return nil
}

func InternalErrorHandler(req *Request, res *Response) error {
	res.StatusCode = 500
	res.SetContent([]byte{})
	return nil
}

type Multiplexer struct {
	word     string
	handler  HandlerFunc
	children map[string]*Multiplexer
}

func NewMultiplexer(handler HandlerFunc) *Multiplexer {
	if handler == nil {
		handler = NotFoundHandler
	}
	return &Multiplexer{
		word:     "/",
		handler:  handler,
		children: map[string]*Multiplexer{},
	}
}

func newMultiplexer(word string, handler HandlerFunc) *Multiplexer {
	return &Multiplexer{
		word:     word,
		handler:  handler,
		children: map[string]*Multiplexer{},
	}
}

func (m *Multiplexer) UpdateHandler(uri string, handler HandlerFunc) {
	if strings.Trim(uri, " ") == "" {
		return
	}
	if strings.Trim(uri, " ") == "/" {
		m.handler = handler
		return
	}
	words := SplitConvertFilter(uri, "/", nil, func(s string) bool { return s != "" })
	curr := m
	for _, v := range words {
		if _, ok := curr.children[v]; !ok {
			curr.children[v] = newMultiplexer(v, nil)
		}
		curr = curr.children[v]
	}
	curr.handler = handler
}

func (m *Multiplexer) findHandler(uri string) *HandlerFunc {
	if strings.Trim(uri, " ") == "/" {
		return &m.handler
	}
	words := SplitConvertFilter(uri, "/", nil, func(s string) bool { return s != "" })
	curr := m
	for _, v := range words {
		if _, ok := curr.children[v]; !ok {
			return &curr.handler
		}
		curr = curr.children[v]
	}
	return &curr.handler
}
