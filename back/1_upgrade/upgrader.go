package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

type WsUpgrader struct{}

func NewWsUpgrader() *WsUpgrader {
	return &WsUpgrader{}
}

func (c *WsUpgrader) getReturnSec(webSecSocketkey string) string {
	var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	h := sha1.New()
	h.Write([]byte(webSecSocketkey))
	h.Write(keyGUID)
	secWebSocketAccept := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return secWebSocketAccept
}

func (c *WsUpgrader) GetResponseString(webSecSocketKey string) []byte {
	response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", webSecSocketKey)
	fmt.Println(response)
	return []byte(response)
}
