package jsonrpc2http

import (
	"net/http"
)

type Server struct {
	*http.Server
	*HTTPHandler

	logger Logger
}

type ServerConfig struct {
	Addr    string
	Handler *HTTPHandler
	Logger  Logger
}

func NewServer(config ServerConfig) *Server {
	if config.Logger == nil {
		config.Logger = new(SimpleLogger)
	}
	if config.Handler == nil {
		config.Handler = NewHTTPHandler(HandlerConfig{
			Logger:     config.Logger,
			HandlerMap: nil,
		})
	}

	return &Server{
		HTTPHandler: config.Handler,
		Server: &http.Server{
			Addr:    config.Addr,
			Handler: config.Handler,
		},
	}
}
