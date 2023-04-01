package main

import "fmt"

var (
	ErrBadRequst     = fmt.Errorf("bad request")
	ErrInternalError = fmt.Errorf("server internal error")
)
