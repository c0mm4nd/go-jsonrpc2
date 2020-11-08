package jsonrpc2http

import (
	"net/http"
)

type Logger interface {
	Debug(...interface{})
	Error(...interface{})
}

type Server struct {
	http.Server
	*HTTPHandler

	logger Logger
}

type ServerConfig struct {
	Addr    string
	Handler *HTTPHandler
	Logger  Logger
}

func NewServer(config ServerConfig) *Server {
	var handler = config.Handler
	if config.Handler == nil {
		handler = NewHTTPHandler(HandlerConfig{
			Logger:     config.Logger,
			HandlerMap: nil,
		})
	}

	return &Server{
		HTTPHandler: handler,
		Server: http.Server{
			Addr:    config.Addr,
			Handler: handler,
		},
	}
}
