package main

import (
	"log"
	"time"

	"github.com/c0mm4nd/go-jsonrpc2/jsonrpc2ws"
	"github.com/gorilla/websocket"

	"github.com/c0mm4nd/go-jsonrpc2"
)

type MyJsonHandler struct{}

func (h *MyJsonHandler) Handle(_ *websocket.Conn, msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	// result, _ := json.Marshal(map[string]interface{}{"ok": true})
	return jsonrpc2.NewJsonRpcSuccess(msg.ID, nil) // never use []byte{}
}

func main() {
	server, _ := jsonrpc2ws.NewServer(jsonrpc2ws.ServerConfig{
		Addr: "127.0.0.1:8888",
	})

	server.RegisterJsonRpcHandler("check", new(MyJsonHandler))
	server.RegisterJsonRpcHandleFunc("checkAgain", func(_ *websocket.Conn, msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
		// result, _ := json.Marshal(map[string]interface{}{"ok": true})
		return jsonrpc2.NewJsonRpcSuccess(msg.ID, nil)
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	client, _ := jsonrpc2ws.NewClient(jsonrpc2ws.ClientConfig{
		Addr: "127.0.0.1:8888",
		Path: "/",
	})

	du := time.Tick(10 * time.Second)
	for {
		<-du
		msgType := websocket.TextMessage

		msg := jsonrpc2.NewJsonRpcRequest(1, "hello", nil)
		err := client.WriteMessage(msgType, msg)
		if err != nil {
			panic(err)
		}

		_, msg, err = client.ReadMessage()
		if err != nil {
			panic(err)
		}

		log.Printf("reply: %#v\n: %v", msg, msg.Error) // error

		msg = jsonrpc2.NewJsonRpcRequest(1, "check", nil)
		err = client.WriteMessage(msgType, msg)
		if err != nil {
			panic(err)
		}

		_, msg, err = client.ReadMessage()
		if err != nil {
			panic(err)
		}

		log.Printf("reply: %#v\n: %v", msg, msg.Error)

		msg = jsonrpc2.NewJsonRpcRequest(1, "checkAgain", nil)
		err = client.WriteMessage(msgType, msg)
		if err != nil {
			panic(err)
		}

		_, msg, err = client.ReadMessage()
		if err != nil {
			panic(err)
		}

		log.Printf("reply: %#v\n: %v", msg, msg.Error)

	}
}
