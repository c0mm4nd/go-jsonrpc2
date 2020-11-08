package jsonrpc2

import (
	"encoding/json"
)

const jsonRpcVersion = "2.0"

type MsgType int

const (
	TypeInvalidMsg MsgType = iota
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

func (m *JsonRpcMessage) GetType() MsgType {
	if m.Version != jsonRpcVersion {
		return TypeInvalidMsg
	}

	if m.Method != "" { //Request
		if m.ID == nil {
			return TypeNotificationMsg
		} else {
			return TypeRequestMsg
		}
	} else {
		if m.Error != nil { // Response
			return TypeErrorMsg
		} else if m.Result != nil {
			return TypeSuccessMsg
		} else {
			return TypeInvalidMsg
		}
	}
}

func UnmarshalMessage(raw []byte) (*JsonRpcMessage, error) {
	var m JsonRpcMessage
	err := json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
