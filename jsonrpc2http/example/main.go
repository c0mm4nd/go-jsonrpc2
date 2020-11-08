package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/c0mm4nd/go-jsonrpc2"
	"github.com/c0mm4nd/go-jsonrpc2/jsonrpc2http"
)

type MyJsonHandler struct {
}

func (h *MyJsonHandler) Handle(msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	//result, _ := json.Marshal(map[string]interface{}{"ok": true})
	return jsonrpc2.NewJsonRpcSuccess(msg.ID, nil) // never use []byte{}
}

func main() {
	server := jsonrpc2http.NewServer(jsonrpc2http.ServerConfig{
		Addr: "127.0.0.1:8888",
	})

	server.RegisterJsonRpcHandler("check", new(MyJsonHandler))
	server.RegisterJsonRpcHandleFunc("checkAgain", func(msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
		//result, _ := json.Marshal(map[string]interface{}{"ok": true})
		return jsonrpc2.NewJsonRpcSuccess(msg.ID, nil)
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	client := jsonrpc2http.NewClient()
	msg := jsonrpc2.NewJsonRpcRequest(1, "check", nil)

	du := time.Tick(10 * time.Second)
	for {
		select {
		case <-du:
			req, err := jsonrpc2http.NewClientRequest("http://127.0.0.1:8888", msg)
			if err != nil {
				panic(err)
			}

			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			log.Println(string(body))
		}
	}
}
