package handmadehttp

import (
	"net"
	"time"

	"golang.org/x/exp/slog"
)

// TODO: howto unit test?
type Server struct {
	network      string
	addr         string
	timeout      time.Duration
	listener     net.Listener
	connChan     chan net.Conn
	mutplexerGet *Multiplexer
	// TODO: add suport to other methods
}

func NewServer(network, addr string, timeout time.Duration) *Server {
	return &Server{
		network:      network,
		addr:         addr,
		timeout:      timeout,
		listener:     nil,
		connChan:     nil,
		mutplexerGet: NewMultiplexer(nil),
	}
}

func (s *Server) AcceptConns(timeout time.Duration) {
	defer func() {
		if e := recover(); e != nil {
			slog.Error("fatal error %s, acceptConns restart", e)
			go s.AcceptConns(timeout)
		}
	}()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			slog.Warn("fail to accept, %s, %s", conn, err)
			conn.Close()
		}
		go func(conn net.Conn) {
			timer := time.NewTimer(timeout)
			select {
			case s.connChan <- conn:
				slog.Debug("enqueue conn %s, chan len %s", conn, len(s.connChan))
			case <-timer.C:
				conn.Close()
				slog.Warn("timeout %s for enqueue chan %s, drop", timeout, conn)
			}
		}(conn)
	}
}

func (s *Server) HandleConns() {

	for {
		conn := <-s.connChan
		go func(conn net.Conn) {
			defer func() {
				if e := recover(); e != nil {
					slog.Error("fatal error %s when handle conn %s", e, conn)
					res := NewResponse(500)
					_, _ = conn.Write(res.ToByte())
				}
			}()
			defer conn.Close()
			request := NewRequest()
			request.Read(conn)
			response := NewResponse(200)
			if request.ReqMethod != ReqMethodGet {
				// 501: not implemented
				response.StatusCode = 501
				slog.Warn("not implemtmented method", request.ReqMethod)
			} else {
				handler := s.mutplexerGet.findHandler(request.URI)
				if err := (*handler)(request, response); err != nil {
					slog.Warn("fail to handle with %s, err %s", *handler, err)
				}

			}
			buff := response.ToByte()
			n, err := conn.Write(buff)
			if err != nil || n != len(buff) {
				slog.Warn("fail to send response, err %s, %d/%d", err, n, len(buff))
			} else {
				slog.Debug("conn %s done.", conn)
			}
		}(conn)
	}
}
func (s *Server) start() error {
	lis, err := net.Listen(s.network, s.addr)
	if err != nil {
		lis.Close()
		return err
	}
	connChan := make(chan net.Conn, MaxChanSize)
	s.listener = lis
	s.connChan = connChan
	return nil
}
func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
	if s.connChan != nil {
		close(s.connChan)
		s.connChan = nil
	}
}
func (s *Server) UpdateHandler(URI string, fn HandlerFunc) {
	s.mutplexerGet.UpdateHandler(URI, fn)
}
func (s *Server) ListenAndServe() error {
	err := s.start()
	defer s.Stop()
	if err != nil {
		return err
	}
	go s.AcceptConns(s.timeout)
	s.HandleConns()
	return nil
}
