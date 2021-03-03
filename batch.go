package jsonrpc2

type JsonRpcMessageBatch []*JsonRpcMessage

// IsBatchMarshal is a helper to check whether the rawbytes are from a jsonrpc message batch
func IsBatchMarshal(raw []byte) bool {
	return raw[0] == '['
}

func NewJsonRpcMessageBatch(messages ...*JsonRpcMessage) JsonRpcMessageBatch {
	return messages
}

func (b *JsonRpcMessageBatch) Marshal() ([]byte, error) {
	return JSON.Marshal(b)
}

func (b *JsonRpcMessageBatch) Unmarshal(raw []byte) (*JsonRpcMessageBatch, error) {
	err := JSON.Unmarshal(raw, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func MarshalMessageBatch(m *JsonRpcMessage) ([]byte, error) {
	return JSON.Marshal(m)
}

func UnmarshalMessageBatch(raw []byte) (JsonRpcMessageBatch, error) {
	var m JsonRpcMessageBatch
	err := JSON.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
