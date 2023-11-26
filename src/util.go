package lib

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func PrintJSON(v any) {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("json marshal error")
	}
	fmt.Println(string(json))
}

func GetUpgradeResponseString(webSecSocketkey string) []byte {
	var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	h := sha1.New()
	h.Write([]byte(webSecSocketkey))
	h.Write(keyGUID)
	secWebSocketAccept := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return []byte(fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", secWebSocketAccept))
}

// func GetCutTail(raw []byte) string {
// 	return regexp.MustCompile(`(.*)\\`).FindStringSubmatch("ggqu123\xaa")[1]
// }
