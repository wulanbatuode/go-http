package handmadehttp

import (
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slog"
)

// TODO: howto unit test?
type Server struct {
	network  string
	addr     string
	timeout  time.Duration
	listener net.Listener
	connChan chan net.Conn
	// TODO: add suport to other methods
	mutplexerGet *Multiplexer
	running      int64
	Wg           *sync.WaitGroup
}

func NewServer(network, addr string, timeout time.Duration) *Server {
	return &Server{
		network:      network,
		addr:         addr,
		timeout:      timeout,
		listener:     nil,
		connChan:     nil,
		mutplexerGet: NewMultiplexer(nil),
		running:      0,
		Wg:           &sync.WaitGroup{},
	}
}

func (s *Server) isRunning() bool {
	return atomic.LoadInt64(&s.running) != 0
}
func (s *Server) setRunningState(target int64) {
	atomic.StoreInt64(&s.running, target)
}
func (s *Server) AcceptConns(timeout time.Duration) {
	defer func() {
		if e := recover(); e != nil {
			slog.Error("fatal error %s, acceptConns restart", e)
			go s.AcceptConns(timeout)
		} else {
			s.Wg.Done()
		}
	}()
	for s.isRunning() {
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

	for s.isRunning() {
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
	s.Wg.Done()
}

// FIXME: race condition
func (s *Server) Stop() {
	slog.Info("stopping server")
	s.setRunningState(0)
	s.Wg.Wait()
	slog.Info("closing listener and chan")
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
	s.mutplexerGet.UpdateHandler(strings.ToUpper(URI), fn)
}
func (s *Server) ListenAndServe() error {
	// TODO: how to notice caller server is ready by chan?
	slog.Info("starting server")
	lis, err := net.Listen(s.network, s.addr)
	if err != nil {
		// lis.Close()
		return err
	}
	connChan := make(chan net.Conn, MaxChanSize)
	s.listener = lis
	s.connChan = connChan
	s.setRunningState(int64(1))
	s.Wg.Add(2)
	go s.AcceptConns(s.timeout)
	go s.HandleConns()
	return nil
}
