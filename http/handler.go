package http

import (
	"github.com/maoxs2/go-jsonrpc2"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler struct {
	handlerMap map[string]JsonRpcHandler
}

func NewHTTPHandler() *Handler {
	return &Handler{handlerMap: map[string]JsonRpcHandler{}}
}

func (h *Handler) RegisterJsonRpcHandleFunc(method string, fn func(*jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage) {
	handler := jsonRpcHandlerFunc(fn)
	h.handlerMap[method] = handler
}

func (h *Handler) RegisterJsonRpcHandler(method string, handler JsonRpcHandler) {
	h.handlerMap[method] = handler
}


func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	raw, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	var res = &jsonrpc2.JsonRpcMessage{}
	jsonRpcReq, err := jsonrpc2.UnmarshalMessage(raw)
	if err != nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrParseFailed, err)
		res = jsonrpc2.NewJsonRpcError(nil, errParams)
	}
	res = h.serveMessage(jsonRpcReq)

	b, err := res.Marshal()
	if err != nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrInternalError, err)
		res = jsonrpc2.NewJsonRpcError(jsonRpcReq.ID, errParams)
	}

	if res.GetType() == jsonrpc2.TypeErrorMsg {
		w.WriteHeader(400)
	}else{
		w.WriteHeader(200)
	}

	w.Write(b)
}

func (h *Handler) serveMessage(req *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	handler, exists := h.handlerMap[req.Method]
	if !exists {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrMethodNotFound)
		return jsonrpc2.NewJsonRpcError(nil, errParams)
	}

	res := handler.Handle(req)
	if res == nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrInternalError)
		return jsonrpc2.NewJsonRpcError(nil, errParams)
	} else {
		return res
	}
}
