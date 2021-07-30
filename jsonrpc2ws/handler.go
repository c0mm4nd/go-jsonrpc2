package jsonrpc2ws

import (
	"encoding/json"
	"net/http"

	"github.com/c0mm4nd/go-jsonrpc2"
	"github.com/gorilla/websocket"
)

// WSHandler is acting as a http.Handler and will redirect the jsonrpc message to one of the registered jsonrpc handlers on its handler table
type WSHandler struct {
	upgrader *websocket.Upgrader
	HandlerConfig
}

type HandlerConfig struct {
	Logger     Logger
	HandlerMap map[string]StatefulJsonRpcHandler
}

func NewWSHandler(config HandlerConfig) *WSHandler {
	logger := config.Logger
	if logger == nil {
		config.Logger = new(SimpleLogger)
	}

	if config.HandlerMap == nil {
		config.HandlerMap = make(map[string]StatefulJsonRpcHandler)
	}

	return &WSHandler{
		upgrader:      &websocket.Upgrader{},
		HandlerConfig: config,
	}
}

func (h *WSHandler) RegisterJsonRpcHandleFunc(method string, fn func(*websocket.Conn, *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage) {
	handler := statefulJsonRpcHandlerFunc(fn)
	h.HandlerMap[method] = handler
}

func (h *WSHandler) RegisterJsonRpcHandler(method string, handler StatefulJsonRpcHandler) {
	h.HandlerMap[method] = handler
}

// ServeHTTP = ServeWS
func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Error(err)
		return
	}
	defer func() {
		err := c.Close()
		if err != nil {
			h.Logger.Error(err)
		}
	}() // serve http like tcp

	for {
		msgType, rawMsg, err := c.ReadMessage()
		if err != nil {
			h.Logger.Error(err)
			break
		}

		if jsonrpc2.IsBatchMarshal(rawMsg) {
			h.onBatchMsg(c, msgType, rawMsg)
		} else {
			h.onSingleMsg(c, msgType, rawMsg)
		}
	}
}

func (h *WSHandler) onSingleMsg(c *websocket.Conn, msgType int, raw []byte) {
	var res *jsonrpc2.JsonRpcMessage
	jsonRPCReq, err := jsonrpc2.UnmarshalMessage(raw)
	if err != nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrParseFailed, err)
		e := jsonrpc2.NewJsonRpcError(nil, errParams)
		b, err := json.Marshal(e)
		if err != nil {
			h.Logger.Error(err)
		}

		err = c.WriteMessage(msgType, b)
		if err != nil {
			h.Logger.Error(err)
		}
	}
	res = h.serveSingleMessage(c, jsonRPCReq)

	b, err := json.Marshal(res)
	if err != nil {
		h.Logger.Error(err)
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrInternalError, err)
		res = jsonrpc2.NewJsonRpcError(nil, errParams)
		b, err = json.Marshal(res)
		if err != nil {
			panic(err)
		}
	}

	err = c.WriteMessage(msgType, b)
	if err != nil {
		h.Logger.Error(err)
	}
}

func (h *WSHandler) onBatchMsg(conn *websocket.Conn, msgType int, raw []byte) {
	res := jsonrpc2.JsonRpcMessageBatch{}
	jsonRPCReqBatch, err := jsonrpc2.UnmarshalMessageBatch(raw)
	if err != nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrParseFailed, err)
		e := jsonrpc2.NewJsonRpcError(nil, errParams)
		b, err := json.Marshal(e)
		if err != nil {
			h.Logger.Error(err)
		}

		err = conn.WriteMessage(msgType, b)
		if err != nil {
			h.Logger.Error(err)
		}
	}
	res = h.serveBatchMessage(conn, jsonRPCReqBatch)

	b, err := res.Marshal()
	if err != nil {
		h.Logger.Error(err)
	}

	err = conn.WriteMessage(msgType, b)
	if err != nil {
		h.Logger.Error(err)
	}
}

func (h *WSHandler) serveSingleMessage(conn *websocket.Conn, req *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	handler, exists := h.HandlerMap[req.Method]
	if !exists {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrMethodNotFound)
		return jsonrpc2.NewJsonRpcError(nil, errParams)
	}

	res := handler.Handle(conn, req)
	if res == nil {
		errParams := jsonrpc2.NewError(0, jsonrpc2.ErrInternalError)
		return jsonrpc2.NewJsonRpcError(nil, errParams)
	}

	return res
}

func (h *WSHandler) serveBatchMessage(conn *websocket.Conn, reqBatch jsonrpc2.JsonRpcMessageBatch) jsonrpc2.JsonRpcMessageBatch {
	resBatch := make(jsonrpc2.JsonRpcMessageBatch, len(reqBatch))
	for i := 0; i < len(reqBatch); i++ {
		handler, exists := h.HandlerMap[reqBatch[i].Method]
		if !exists {
			errParams := jsonrpc2.NewError(0, jsonrpc2.ErrMethodNotFound)
			resBatch[i] = jsonrpc2.NewJsonRpcError(nil, errParams)
			continue
		}

		res := handler.Handle(conn, reqBatch[i])
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
