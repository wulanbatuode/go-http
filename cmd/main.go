package main

import (
	"fmt"
	"handmadehttp/pkg/handmadehttp"
	"time"

	"golang.org/x/exp/slog"
)

func main() {
	slog.Info("start server")
	s := handmadehttp.NewServer("tcp", ":8080", 5*time.Second)
	s.UpdateHandler("/echo",
		func(req *handmadehttp.Request, res *handmadehttp.Response) error {
			res.SetContent([]byte(fmt.Sprintf("%s %s", req.URI, req.Param)))
			return nil
		})
	go s.ListenAndServe()
	// if err != nil {
	// 	slog.Error("fail to start server %s", err)
	// }
	slog.Info("xxxxx sleep 1 sec")
	time.Sleep(time.Second)
	slog.Info("stopping")
	s.Stop()
	slog.Info("wating")
	s.Wg.Wait()
	slog.Info("finish")
}
