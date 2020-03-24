package jsonrpc2http

import "github.com/maoxs2/go-jsonrpc2"

type JsonRpcHandler interface {
	Handle(*jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage
}

type jsonRpcHandlerFunc func(*jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage

func (fn jsonRpcHandlerFunc) Handle(msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	return fn(msg)
}
