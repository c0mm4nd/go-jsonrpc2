package jsonrpc2

import (
	"encoding/json"
	"reflect"
)

// WithParams converts the params into raw bytes and store in the message
func (msg *JsonRpcMessage) WithParams(params interface{}) error {
	rawParams, err := JSON.Marshal(params)
	if err != nil {
		return err
	}

	paramsField := json.RawMessage(rawParams)
	msg.Params = &paramsField

	return nil
}

// LoadParams reads the params from raw bytes to the pointer
func (msg *JsonRpcMessage) LoadParams(paramsPtr interface{}) error {
	if reflect.ValueOf(paramsPtr).Kind() != reflect.Ptr {
		panic("the params should be a pointer in JsonRpcMessage.LoadParams")
	}

	err := JSON.Unmarshal(*msg.Params, paramsPtr)
	if err != nil {
		return err
	}

	return nil
}
