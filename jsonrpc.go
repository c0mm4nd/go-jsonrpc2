package jsonrpc2

import json "encoding/json"

const jsonRpcVersion = "2.0"

type JsonMsgType int

const (
	TypeInvalidMsg JsonMsgType = iota
	TypeNotificationMsg
	TypeRequestMsg
	TypeErrorMsg
	TypeSuccessMsg
)

type JsonRpcMessage struct {
	Version string `json:"jsonrpc"`

	Method string `json:"method,omitempty"`

	Params *json.RawMessage `json:"params,omitempty"`
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *Error           `json:"error,omitempty"`

	ID interface{} `json:"id,omitempty"`
}

func (m *JsonRpcMessage) GetType() JsonMsgType {
	if m.Version != jsonRpcVersion {
		return TypeInvalidMsg
	}

	if m.Method != "" { //Request
		if m.ID == nil {
			return TypeNotificationMsg
		}

		return TypeRequestMsg
	} else {
		if m.Error != nil { // Response
			return TypeErrorMsg
		} else if m.Result != nil {
			return TypeSuccessMsg
		}
		return TypeInvalidMsg
	}
}

func UnmarshalMessage(raw []byte) (*JsonRpcMessage, error) {
	var m JsonRpcMessage
	err := JSON.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
