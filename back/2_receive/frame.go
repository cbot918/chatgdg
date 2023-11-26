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
	c.PayloadLen = secondByte & 0b01111111

	// process mask
	mask := []byte{data[2], data[3], data[4], data[5]}
	fmt.Println("mask: ", mask)

	// process payload data
	payload := []byte{}
	for i := 6; i <= int(c.PayloadLen+6); i++ {
		payload = append(payload, data[i])
	}
	fmt.Println("payload: ", payload)

	// XOR payload and mask
	result := []byte{}
	for i, item := range payload {
		result = append(result, item^mask[i%4])
	}

	return result
}
