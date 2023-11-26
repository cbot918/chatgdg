package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

const network = "tcp"
const port = ":8890"

var count int

func main() {

	fmt.Println(port)
	listener, err := net.Listen(network, port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {

			/* HTTP to WS */
			data, err := read(conn)
			if err != nil {
				fmt.Println("connect failed")
				return
			}

			key := getWebSecKey(data)
			retKey := getReturnSec(key)
			response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", retKey)

			_, err = conn.Write([]byte(response))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("connect established")
			/* HTTP to WS */

			for {
				// /* WS Read Write */
				data, err = read(conn)
				if err != nil {
					continue
				}
				fmt.Printf("%q", string(decodeFrame(data)))
			}

		}(conn)
	}

}

func read(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, err
}

func getWebSecKey(data []byte) string {
	pattern := `Sec-WebSocket-Key: ([^\r\n]+)`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(string(data))
	return strings.TrimSpace(match[1])
}

func getReturnSec(webSecSocketkey string) string {
	var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	h := sha1.New()
	h.Write([]byte(webSecSocketkey))
	h.Write(keyGUID)
	secWebSocketAccept := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return secWebSocketAccept
}

type Frame struct {
	Fin        byte
	Opcode     byte
	IsMasked   byte
	PayloadLen byte
	Mask       []byte
	Payload    []byte
}

func decodeFrame(data []byte) []byte {
	firstByte := data[0]
	secondByte := data[1]

	f := new(Frame)

	f.Fin = firstByte & 0b10000000
	f.Opcode = firstByte & 0b00001111
	f.IsMasked = secondByte & 0b10000000
	f.PayloadLen = secondByte & 0b01111111

	// process mask
	mask := []byte{data[2], data[3], data[4], data[5]}
	// fmt.Println("mask: ", mask)

	// process payload data
	payload := []byte{}
	for i := 6; i <= int(f.PayloadLen+6); i++ {
		payload = append(payload, data[i])
	}
	// fmt.Println("payload: ", payload)

	// XOR payload and mask
	result := []byte{}
	for i, item := range payload {
		result = append(result, item^mask[i%4])
	}

	return result
}
