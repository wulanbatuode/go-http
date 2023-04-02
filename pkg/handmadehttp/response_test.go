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
	assert.Equal(t, "HTTP/1.1 500 Internal Server Error\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\nContent-Length: 0\r\n\r\n",
		string(res.ToByte()))
	res.AppendContent([]byte("hello"))
	res.StatusCode = 200
	assert.Equal(t, len("hello"), res.ContentLength)
	assert.Equal(t, []byte("hello"), *res.Body)
	expect := "HTTP/1.1 200 OK\r\nContent-Type: text/html; charset=UTF-8\r\n" +
		"Content-Length: 5\r\n\r\nhello"
	assert.Equal(t, expect, string(res.ToByte()))
}
