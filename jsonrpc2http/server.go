package jsonrpc2http

import "net/http"

type Server struct {
	http.Server
}

func NewServer(addr string, handler *HTTPHandler) *Server {
	return &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}
