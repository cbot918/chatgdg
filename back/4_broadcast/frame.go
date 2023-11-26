package main

import "fmt"

type Frame struct {
	Fin        byte
	Opcode     byte
	IsMasked   byte
	PayloadLen byte
	Mask       []byte
	Payload    []byte
}

func NewFrame() *Frame {
	return &Frame{}
}

func (c *Frame) DecodeFrame(data []byte) []byte {
	firstByte := data[0]
	secondByte := data[1]

	c.Fin = firstByte & 0b10000000
	c.Opcode = firstByte & 0b00001111
	c.IsMasked = secondByte & 0b10000000
	payloadLength := int(secondByte & 0b01111111)

	var maskStart int
	if payloadLength == 126 {
		// Next 2 bytes are the length
		payloadLength = int(data[2])<<8 | int(data[3])
		maskStart = 4
	} else if payloadLength == 127 {
		// Next 8 bytes are the length
		// Assuming payload is small enough for int
		payloadLength = 0
		for i := 0; i < 8; i++ {
			payloadLength = payloadLength<<8 | int(data[2+i])
		}
		maskStart = 10
	} else {
		maskStart = 2
	}

	// Ensure there is enough data for the mask and payload
	if len(data) < maskStart+4+payloadLength {
		fmt.Println("Error: Data too short for expected frame size")
		return nil
	}

	// process mask
	mask := data[maskStart : maskStart+4]

	// process payload data
	payloadStart := maskStart + 4
	payload := data[payloadStart : payloadStart+payloadLength]

	// XOR payload and mask
	result := make([]byte, payloadLength)
	for i := 0; i < payloadLength; i++ {
		result[i] = payload[i] ^ mask[i%4]
	}

	return result
}

func (c *Frame) EncodeFrame(message []byte) []byte {
	// Create a buffer to hold the WebSocket frame
	frame := []byte{0x81} // Text frame (FIN bit set)

	// Calculate the payload length
	payloadLen := len(message)

	// Add the payload length to the frame
	if payloadLen <= 125 {
		frame = append(frame, byte(payloadLen))
	} else if payloadLen <= 65535 {
		frame = append(frame, 126) // 16-bit payload length
		frame = append(frame, byte(payloadLen>>8), byte(payloadLen&0xFF))
	} else {
		frame = append(frame, 127) // 64-bit payload length
		for i := 0; i < 8; i++ {
			frame = append(frame, byte(payloadLen>>(56-(i*8))))
		}
	}

	// Add the message data to the frame
	frame = append(frame, message...)

	return frame
}

// func (c *Frame) DecodeFrame(data []byte) []byte {
// 	firstByte := data[0]
// 	secondByte := data[1]

// 	c.Fin = firstByte & 0b10000000
// 	c.Opcode = firstByte & 0b00001111
// 	c.IsMasked = secondByte & 0b10000000
// 	c.PayloadLen = secondByte & 0b01111111

// 	// process mask
// 	mask := []byte{data[2], data[3], data[4], data[5]}
// 	fmt.Println("mask: ", mask)

// 	// process payload data
// 	payload := []byte{}
// 	for i := 6; i <= int(c.PayloadLen+6); i++ {
// 		payload = append(payload, data[i])
// 	}
// 	fmt.Println("payload: ", payload)

// 	// XOR payload and mask
// 	result := []byte{}
// 	for i, item := range payload {
// 		result = append(result, item^mask[i%4])
// 	}

// 	return result
// }
