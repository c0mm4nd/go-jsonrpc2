package jsonrpc2ws

import (
	"net/http"
)

type Server struct {
	*http.Server
	*WSHandler

	logger Logger
}

type ServerConfig struct {
	Addr string
	Path string

	Handler *WSHandler
	Logger  Logger
}

func NewServer(config ServerConfig) (*Server, error) {
	if config.Logger == nil {
		config.Logger = new(SimpleLogger)
	}
	if config.Handler == nil {
		config.Handler = NewWSHandler(HandlerConfig{
			Logger:     config.Logger,
			HandlerMap: nil,
		})
	}

	server := &Server{
		WSHandler: config.Handler,
		logger:    config.Logger,
		Server: &http.Server{
			Addr:    config.Addr,
			Handler: config.Handler,
		},
	}

	return server, nil
}

func WrapHTTP(httpServer *http.Server, config ServerConfig) (*Server, error) {
	if config.Logger == nil {
		config.Logger = new(SimpleLogger)
	}
	if config.Handler == nil {
		config.Handler = NewWSHandler(HandlerConfig{
			Logger:     config.Logger,
			HandlerMap: nil,
		})
	}

	server := &Server{
		WSHandler: config.Handler,
		logger:    config.Logger,
		Server:    httpServer,
	}

	return server, nil
}
