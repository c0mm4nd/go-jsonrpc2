package jsonrpc2

import (
	"encoding/json"
)

// NewJsonRpcRequest returns a JSON-RPC 2.0 request message structures. id must be string/int/nil type. params should be json marshaled
func NewJsonRpcRequest(id interface{}, method string, params []byte) *JsonRpcMessage {
	paramsField := json.RawMessage(params)

	p := &JsonRpcMessage{
		Version: jsonRpcVersion,
		Method:  method,
		Params:  &paramsField,
		ID:      id,
	}

	return p
}

// NewJsonRpcNotification returns a JSON-RPC 2.0 notification message structures which doesnt have id. params should be json marshaled
func NewJsonRpcNotification(method string, params []byte) *JsonRpcMessage {
	return NewJsonRpcRequest(nil, method, params)
}

// NewJsonRpcSuccess returns a JSON-RPC 2.0 success message structures. result should be json marshaled
func NewJsonRpcSuccess(id interface{}, result []byte) *JsonRpcMessage {
	resultField := json.RawMessage(result)

	p := &JsonRpcMessage{
		Version: jsonRpcVersion,
		Result:  &resultField,
		ID:      id,
	}

	return p
}

// NewJsonRpcError returns a JSON-RPC 2.0 error message structures.
func NewJsonRpcError(id interface{}, errParams *Error) *JsonRpcMessage {
	p := &JsonRpcMessage{
		Version: jsonRpcVersion,
		Error:   errParams,
		ID:      id,
	}

	return p
}
