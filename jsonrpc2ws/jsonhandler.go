package jsonrpc2ws

import (
	"github.com/c0mm4nd/go-jsonrpc2"
	"github.com/gorilla/websocket"
)

type StatefulJsonRpcHandler interface {
	Handle(*websocket.Conn, *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage
}

type statefulJsonRpcHandlerFunc func(*websocket.Conn, *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage

func (fn statefulJsonRpcHandlerFunc) Handle(conn *websocket.Conn, msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	return fn(conn, msg)
}
