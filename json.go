package jsonrpc2

import json "encoding/json"

type stdJSON struct{}

func (stdJSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (stdJSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// JSON is the json parser inside the jsonrpc2,
// developers can change parser via assign another (un)marshaller into this symbol
var JSON = new(stdJSON)
