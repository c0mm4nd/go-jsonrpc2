package main

import (
	"encoding/json"
	"github.com/maoxs2/go-jsonrpc2"
	"github.com/maoxs2/go-jsonrpc2/http"
	"time"
)

type MyJsonHandler struct {
}

func (h *MyJsonHandler) Handle(msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	result, _ := json.Marshal(map[string]interface{}{"ok": true})
	return jsonrpc2.NewJsonRpcSuccess(msg.ID, result)
}

func main() {
	handler := http.NewHTTPHandler()
	handler.RegisterJsonRpcHandler("check", new(MyJsonHandler))
	handler.RegisterJsonRpcHandleFunc("checkAgain",
		func(msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
			result, _ := json.Marshal(map[string]interface{}{"ok": true})
			return jsonrpc2.NewJsonRpcSuccess(msg.ID, result)
		})

	server := http.NewServer("127:0.0.1:8888", handler)

	go server.ListenAndServe()

	client := http.NewClient()
	msg := jsonrpc2.NewJsonRpcRequest(1, "check", jsonrpc2.EmptyArrayBytes)
	req, _ := http.NewClientRequest("127.0.0.1:8888", msg)

	du := time.Tick(10*time.Second)
	for {
		select {
		case <-du:
			client.Do(req)
		}
	}
}
