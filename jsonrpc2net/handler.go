package jsonrpc2net

import (
	"github.com/maoxs2/go-jsonrpc2"
	"io"
	"log"
)

type JsonRpcHandler interface {
	Handle(*jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage
}

type jsonRpcHandlerFunc func(*jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage

func (fn jsonRpcHandlerFunc) Handle(msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	return fn(msg)
}

func (s *Server) RegisterJsonRpcHandleFunc(method string, fn func(*jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage) {
	handler := jsonRpcHandlerFunc(fn)
	s.handlerMap[method] = handler
}

func (s *Server) RegisterJsonRpcHandler(method string, handler JsonRpcHandler) {
	s.handlerMap[method] = handler
}

func (s *Server) onBatchMsg(w io.Writer, raw []byte) {
	var res = jsonrpc2.JsonRpcMessageBatch{}
	jsonRpcReqBatch, err := jsonrpc2.UnmarshalMessageBatch(raw)
	if err != nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrParseFailed, err)
		e := jsonrpc2.NewJsonRpcError(nil, errParams)
		b, err := e.Marshal()
		if err != nil {
			log.Println(err)
		}
		w.Write(b)
	}
	res = s.serveBatchMessage(jsonRpcReqBatch)

	b, err := res.Marshal()
	if err != nil {
		log.Println(err)
	}

	w.Write(b)
}

func (s *Server) onSingleMsg(w io.Writer, raw []byte) {
	var res = &jsonrpc2.JsonRpcMessage{}
	jsonRpcReq, err := jsonrpc2.UnmarshalMessage(raw)
	if err != nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrParseFailed, err)
		res = jsonrpc2.NewJsonRpcError(nil, errParams)
	}
	res = s.serveSingleMessage(jsonRpcReq)

	b, err := res.Marshal()
	if err != nil {
		log.Println(err)
	}

	w.Write(b)
}

func (s *Server) serveSingleMessage(req *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	handler, exists := s.handlerMap[req.Method]
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

func (s *Server) serveBatchMessage(reqBatch jsonrpc2.JsonRpcMessageBatch) jsonrpc2.JsonRpcMessageBatch {
	var resBatch = make(jsonrpc2.JsonRpcMessageBatch, len(reqBatch))
	for i := 0; i < len(reqBatch); i++ {
		handler, exists := s.handlerMap[reqBatch[i].Method]
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
