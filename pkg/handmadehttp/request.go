package handmadehttp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	Method        string
	URI           string
	Protocal      string
	ContentType   string
	ContentLength int
	Attr          map[string]interface{}
	Body          *[]byte
}

func NewRequest() *Request {
	return &Request{
		Method:        "",
		URI:           "",
		Protocal:      DefaultProtocal,
		ContentType:   DefaultContentType,
		ContentLength: 0,
		Attr:          map[string]interface{}{},
		Body:          &[]byte{},
	}
}
func (req *Request) readHeader(rd io.Reader) (err error) {
	reader := bufio.NewReader(rd)
	buff := make([]byte, 0, BuffSize)
	for i := 0; i < MaxHeaderLine; i++ {
		line, prefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		if len(buff) == 0 && len(line) == 0 && !prefix {
			// err = nil
			break
		}
		buff = append(buff, line...)
		if prefix {
			continue
		}
		if req.URI == "" {
			tokens := SplitConvertFilter(strings.ToLower(string(buff)), " ", nil, func(s string) bool { return s != "" })
			req.Method = tokens[0]
			req.URI = tokens[1]
			if len(tokens) > 2 {
				req.Protocal = tokens[2]
			}
			buff = buff[:0]
		}
		if len(buff) > 0 {
			tokens := SplitConvertFilter(strings.ToLower(string(buff)), ":", nil, func(s string) bool { return s != "" })
			if len(tokens) != 2 {
				err = fmt.Errorf("%w with attr %s", ErrBadRequst, tokens)
				return err
			}
			switch tokens[0] {
			case KeyContentLength:
				length, err := strconv.Atoi(strings.Trim(tokens[1], " "))
				if err != nil {
					break
				}
				req.ContentLength = length
			case KeyContentType:
				req.ContentType = strings.Trim(tokens[1], " ")
			default:
				req.Attr[strings.Trim(tokens[0], " ")] = strings.Trim(tokens[1], " ")
			}
			buff = buff[:0]
		}
	}
	return err
}

func (req *Request) readBody(rd io.Reader) (err error) {
	buff := make([]byte, req.ContentLength)
	n, err := rd.Read(buff)
	if err != nil && err != io.EOF {
		return err
	}
	if n != req.ContentLength {
		err = fmt.Errorf("%w, expect %d length, got %d", ErrBadRequst, req.ContentLength, n)
		return err
	}
	req.Body = &buff
	return nil
}

func (req *Request) Read(rd io.Reader) (err error) {
	reader := bufio.NewReader(rd)
	err = req.readHeader(reader)
	if err != nil {
		return err
	}
	err = req.readBody(reader)

	if err != nil {
		return err
	}
	return nil
}
