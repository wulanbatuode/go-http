package handmadehttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	res := NewResponse(500)
	assert.Equal(t, 500, res.StatusCode)
	assert.Equal(t, 0, res.ContentLength)
	assert.Equal(t, []byte{}, *res.Body)
	assert.Equal(t, transformCase("HTTP/1.1 500 Internal Server Error\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\nContent-Length: 0\r\n\r\n"),
		string(res.ToByte()))
	res.AppendContent([]byte("Hello"))
	res.StatusCode = 200
	assert.Equal(t, len("Hello"), res.ContentLength)
	assert.Equal(t, []byte("Hello"), *res.Body)
	expect := transformCase("HTTP/1.1 200 OK\r\nContent-Type: text/html; charset=UTF-8\r\n"+
		"Content-Length: 5\r\n\r\n") + "Hello"
	assert.Equal(t, expect, string(res.ToByte()))
}
