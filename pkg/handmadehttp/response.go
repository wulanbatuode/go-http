package handmadehttp

import (
	"fmt"
	"strings"
)

type Response struct {
	ContentType   string
	StatusCode    int
	ContentLength int
	Protocal      string
	Attr          map[string]interface{}
	Body          *[]byte
}

func NewResponse(status int) *Response {
	// TODO: check if status exists in CodeToStatus
	return &Response{
		StatusCode:    status,
		ContentLength: 0,
		ContentType:   DefaultContentType,
		Protocal:      DefaultProtocal,
		Attr:          map[string]interface{}{},
		Body:          &[]byte{},
	}
}
func (res *Response) AppendContent(b []byte) {
	(*res.Body) = append((*res.Body), b...)
	res.ContentLength += len(b)
}
func (res *Response) SetContent(b []byte) {
	res.Body = &b
	res.ContentLength = len(b)
}

func (res *Response) ToByte() []byte {
	header := make([]string, 0, len(res.Attr)+2)
	header = append(header, fmt.Sprintf("%s %d %s%s"+"%s: %s%s"+"%s: %d%s",
		res.Protocal, res.StatusCode, CodeToStatus[res.StatusCode], LineBreak,
		KeyContentType, res.ContentType, LineBreak,
		KeyContentLength, res.ContentLength, LineBreak))
	for k, v := range res.Attr {
		header = append(header, fmt.Sprintf("%s: %d%s", k, v, LineBreak))
	}
	header = append(header, LineBreak)
	buff := make([]byte, 0, len(header)+res.ContentLength)
	buff = append(buff, []byte((strings.Join(header, "")))...)
	buff = append(buff, *res.Body...)
	return buff
}
