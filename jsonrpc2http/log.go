package jsonrpc2http

import "log"

type Logger interface {
	Debug(...interface{})
	Error(...interface{})
}

type SimpleLogger struct{}

func (logger *SimpleLogger) Debug(args ...interface{}) {
	log.Println(args...)
}

func (logger *SimpleLogger) Error(args ...interface{}) {
	log.Println(args...)
}
