package util

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
)

type WsUpgrader struct{}

func NewWsUpgrader() *WsUpgrader {
	return &WsUpgrader{}
}

func (c *WsUpgrader) getReturnSec(webSecSocketkey string) string {
	// secWebSocketKey := strings.TrimSpace(strings.Join(webSecSocketkey, ""))
	var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	h := sha1.New()
	h.Write([]byte(webSecSocketkey))
	h.Write(keyGUID)
	secWebSocketAccept := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return secWebSocketAccept
}

func (c *WsUpgrader) Upgrade(conn net.Conn, webSecSocketKey string) error {
	secWebSocketAccept := c.getReturnSec(webSecSocketKey)
	response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", secWebSocketAccept)
	fmt.Println(response)
	res := []byte(response)
	_, err := conn.Write(res)
	return err
}
