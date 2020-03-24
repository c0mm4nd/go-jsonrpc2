package jsonrpc2http

import (
	"bytes"
	"github.com/maoxs2/go-jsonrpc2"
	"net/http"
)

type Client struct {
	http.Client
}

func NewClientRequest(url string, message *jsonrpc2.JsonRpcMessage) (*http.Request, error) {
	raw, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", url, bytes.NewReader(raw))
}

func NewClient() *Client {
	return &Client{
		Client: http.Client{},
	}
}
