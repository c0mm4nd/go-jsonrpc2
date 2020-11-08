package jsonrpc2

import "encoding/json"

type JsonRpcMessageBatch []*JsonRpcMessage

// IsBatchMarshal is a helper to check whether the rawbytes are from a jsonrpc message batch
func IsBatchMarshal(raw []byte) bool {
	return raw[0] == '['
}

func NewJsonRpcMessageBatch(messages ...*JsonRpcMessage) JsonRpcMessageBatch {
	return messages
}

func (b *JsonRpcMessageBatch) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

func (b *JsonRpcMessageBatch) Unmarshal(raw []byte) (*JsonRpcMessageBatch, error) {
	err := json.Unmarshal(raw, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func MarshalMessageBatch(m *JsonRpcMessage) ([]byte, error) {
	return json.Marshal(m)
}

func UnmarshalMessageBatch(raw []byte) (JsonRpcMessageBatch, error) {
	var m JsonRpcMessageBatch
	err := json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
