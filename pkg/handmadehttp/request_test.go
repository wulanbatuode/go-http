package handmadehttp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestRead(t *testing.T) {
	input := []byte("GET /TEST HTTP/1.1\r\nContent-Type: text/html; charset=UTF-8\r\n" +
		"Content-Length: 13\r\n" +
		"User-Agent: curl/7.16.3 libcurl/7.16.3 OpenSSL/0.9.7l zlib/1.2.3\r\n" +
		"Accept-Language: en, mi\r\n" +
		"\r\n" +
		"Hello, World!")

	rd := bytes.NewReader(input)
	req := NewRequest()
	err := req.Read(rd)
	assert.Equal(t, nil, err)
	assert.Equal(t, 13, req.ContentLength)
	assert.Equal(t, DefaultContentType, req.ContentType)
	assert.Equal(t, "HTTP/1.1", req.Protocal)
	assert.Equal(t, "curl/7.16.3 libcurl/7.16.3 OpenSSL/0.9.7l zlib/1.2.3", req.Attr["User-Agent"])
	assert.Equal(t, "en, mi", req.Attr["Accept-Language"])
	assert.Equal(t, "Hello, World!", string(*req.Body))

}
