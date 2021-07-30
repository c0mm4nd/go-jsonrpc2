package jsonrpc2ws

import (
	"encoding/json"
	"net/url"

	"github.com/c0mm4nd/go-jsonrpc2"
	"github.com/gorilla/websocket"
)

type ClientConfig struct {
	Addr string
	Path string
}

type Client struct {
	*websocket.Conn
	ClientConfig
}

func NewClient(config ClientConfig) (*Client, error) {
	u := url.URL{Scheme: "ws", Host: config.Addr, Path: config.Path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		Conn:         c,
		ClientConfig: config,
	}, nil
}

func (c *Client) ReadRawMessage() (messageType int, message []byte, err error) {
	return c.Conn.ReadMessage()
}

func (c *Client) ReadMessage() (messageType int, message *jsonrpc2.JsonRpcMessage, err error) {
	messageType, rawMsg, err := c.Conn.ReadMessage()
	if err != nil {
		return
	}

	var msg jsonrpc2.JsonRpcMessage
	err = json.Unmarshal(rawMsg, &msg)
	if err != nil {
		return messageType, nil, err
	}

	return messageType, &msg, nil
}

func (c *Client) ReadMessageBatch() (messageType int, message *jsonrpc2.JsonRpcMessageBatch, err error) {
	messageType, rawMsgBatch, err := c.Conn.ReadMessage()
	if err != nil {
		return
	}

	var msgBatch jsonrpc2.JsonRpcMessageBatch
	err = json.Unmarshal(rawMsgBatch, &msgBatch)
	if err != nil {
		return messageType, nil, err
	}

	return messageType, &msgBatch, nil
}

func (c *Client) WriteMessage(messageType int, message *jsonrpc2.JsonRpcMessage) error {
	raw, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = c.Conn.WriteMessage(messageType, raw)
	if err != nil {
		return err
	}

	return nil
}
