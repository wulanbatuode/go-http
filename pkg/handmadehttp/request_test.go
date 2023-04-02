package handmadehttp

import (
	"bytes"
	"strings"
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
	assert.Equal(t, transformCase(DefaultContentType), req.ContentType)
	assert.Equal(t, transformCase("http/1.1"), req.Protocal)
	assert.Equal(t, transformCase("curl/7.16.3 libcurl/7.16.3 OpenSSL/0.9.7l zlib/1.2.3"), req.Attr[transformCase("user-agent")])
	assert.Equal(t, transformCase("en, mi"), req.Attr[transformCase("accept-language")])
	assert.Equal(t, "Hello, World!", string(*req.Body))

}
func transformCase(s string) string {
	return strings.ToUpper(s)
}
