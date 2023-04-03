package handmadehttp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func echoServer(network, addr string) *Server {
	s := NewServer(network, addr, 5*time.Second)
	s.UpdateHandler("/echo",
		func(req *Request, res *Response) error {
			res.SetContent([]byte(fmt.Sprintf("%s %s", req.URI, req.Param)))
			return nil
		})
	return s
}

func TestGetWithParam(t *testing.T) {
	s := echoServer("tcp", "localhost:8080")
	assert.NotNil(t, s)
	go s.ListenAndServe()
	t.Log("server start1")
	//wait for server to start up
	time.Sleep(1 * time.Second)
	resp, err := http.Post("http://localhost:8080/echo",
		"application/x-www-form-urlencoded", bytes.NewReader([]byte("a=4")))
	assert.Nil(t, err)
	assert.Equal(t, 501, resp.StatusCode)
	resp, err = http.Get("http://localhost:8080/")
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	buff := make([]byte, BuffSize)
	for i := 0; i < 10; i++ {
		// resp, err := http.Get(fmt.Sprintf("%s%d", "http://localhost:8080/echo?key=", i))
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/echo?key=%d", i))
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		n, err := resp.Body.Read(buff)

		assert.Equal(t, io.EOF, err)
		assert.Equal(t, fmt.Sprintf("/ECHO map[KEY:%d]", i), string(buff[:n]))
		resp.Body.Close()
	}
	t.Cleanup(func() {
		if s != nil {
			s.Stop()
		}
	})
}

var onceHmt sync.Once

func BenchmarkEchoServer(b *testing.B) {
	var s *Server
	onceHmt.Do(
		func() {
			s = echoServer("tcp", "localhost:8080")
			if s == nil {
				b.FailNow()
			}
			go s.ListenAndServe()
			//wait for server to start up
			time.Sleep(time.Second)
		},
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = http.Get(fmt.Sprintf("%s%d", "http://localhost:8080/echo?key=", i))
		// assert.Equal()
	}
	b.Cleanup(func() {
		if s != nil {
			s.Stop()
		}
	})
}

var onceDefault sync.Once

func BenchmarkDefaultEchoServer(b *testing.B) {
	onceDefault.Do(
		func() {
			http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
				queryValues, err := url.ParseQuery(r.URL.RawQuery)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				key := queryValues.Get("key")
				fmt.Fprintf(w, "key: %s\n", key)
			})
			fmt.Println("default http starting")
			go http.ListenAndServe("localhost:8081", nil)
		},
	)
	fmt.Println("reset timer")

	//wait for server to start up
	time.Sleep(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = http.Get(fmt.Sprintf("%s%d", "http://localhost:8081/echo?key=", i))
	}
}
