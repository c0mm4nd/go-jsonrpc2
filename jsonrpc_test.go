package jsonrpc2

import (
	"encoding/json"
	"testing"
)

func TestMarshal(t *testing.T) {
	res := NewJsonRpcSuccess(1, nil)
	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	t.Log(string(b))

	res = NewJsonRpcSuccess(1, b)
	b, err = json.Marshal(res)
	if err != nil {
		panic(err)
	}
	t.Log(string(b))
}

func TestUnmarshal(t *testing.T) {
	raw := `{"jsonrpc":"2.0","result":null,"id":1}`
	var res JsonRpcMessage
	err := json.Unmarshal([]byte(raw), &res)
	if err != nil {
		panic(err)
	}
	t.Log(res)
}
