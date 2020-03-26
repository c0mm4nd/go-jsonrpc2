package jsonrpc2http

import (
	"github.com/maoxs2/go-jsonrpc2"
	"io/ioutil"
	"log"
	"net/http"
)

// HTTPHandler is acting as a http.Handler and will redirect the jsonrpc message to one of the registered jsonrpc handlers on its handler table
type HTTPHandler struct {
	handlerMap map[string]JsonRpcHandler
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{
		handlerMap: map[string]JsonRpcHandler{},
	}
}

func (h *HTTPHandler) RegisterJsonRpcHandleFunc(method string, fn func(*jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage) {
	handler := jsonRpcHandlerFunc(fn)
	h.handlerMap[method] = handler
}

func (h *HTTPHandler) RegisterJsonRpcHandler(method string, handler JsonRpcHandler) {
	h.handlerMap[method] = handler
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	raw, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if jsonrpc2.IsBatchMarshal(raw) {
		h.onBatchMsg(w, raw)
	} else {
		h.onSingleMsg(w, raw)
	}
}

func (h *HTTPHandler) onSingleMsg(w http.ResponseWriter, raw []byte) {
	var res = &jsonrpc2.JsonRpcMessage{}
	jsonRpcReq, err := jsonrpc2.UnmarshalMessage(raw)
	if err != nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrParseFailed, err)
		res = jsonrpc2.NewJsonRpcError(nil, errParams)
	}
	res = h.serveSingleMessage(jsonRpcReq)

	b, err := res.Marshal()
	if err != nil {
		log.Println(err)
	}

	if res.GetType() == jsonrpc2.TypeErrorMsg {
		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
	}

	w.Write(b)
}

func (h *HTTPHandler) onBatchMsg(w http.ResponseWriter, raw []byte) {
	var res = jsonrpc2.JsonRpcMessageBatch{}
	jsonRpcReqBatch, err := jsonrpc2.UnmarshalMessageBatch(raw)
	if err != nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrParseFailed, err)
		e := jsonrpc2.NewJsonRpcError(nil, errParams)
		b, err := e.Marshal()
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(400)
		w.Write(b)
	}
	res = h.serveBatchMessage(jsonRpcReqBatch)

	b, err := res.Marshal()
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(200)
	w.Write(b)
}

func (h *HTTPHandler) serveSingleMessage(req *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
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

func (h *HTTPHandler) serveBatchMessage(reqBatch jsonrpc2.JsonRpcMessageBatch) jsonrpc2.JsonRpcMessageBatch {
	var resBatch = make(jsonrpc2.JsonRpcMessageBatch, len(reqBatch))
	for i := 0; i < len(reqBatch); i++ {
		handler, exists := h.handlerMap[reqBatch[i].Method]
		if !exists {
			errParams := jsonrpc2.NewError(0, jsonrpc2.ErrMethodNotFound)
			resBatch[i] = jsonrpc2.NewJsonRpcError(nil, errParams)
			continue
		}

		res := handler.Handle(reqBatch[i])
		if res == nil {
			errParams := jsonrpc2.NewError(0, jsonrpc2.ErrInternalError)
			resBatch[i] = jsonrpc2.NewJsonRpcError(nil, errParams)
			continue
		} else {
			resBatch[i] = res
			continue
		}
	}

	return resBatch
}
