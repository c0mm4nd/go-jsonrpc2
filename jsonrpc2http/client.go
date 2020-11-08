package jsonrpc2http

import (
	"bytes"
	"net/http"

	"github.com/c0mm4nd/go-jsonrpc2"
)

type Client struct {
	http.Client
}

func NewClientRequest(url string, message *jsonrpc2.JsonRpcMessage) (*http.Request, error) {
	raw, err := message.MarshalJSON()
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
