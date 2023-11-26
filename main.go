package main

import (
	"bufio"
	lib "chatgdg/src"
	"fmt"
	"log"
	"net"
)

const (
	tcp    = "tcp"
	domain = ":8888"
)

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan []byte)
	join      = make(chan *Client)
	leave     = make(chan *Client)
)

type Client struct {
	conn   net.Conn
	writer *bufio.Writer
}

func main() {

	fmt.Println("listening:", domain)
	lis, err := net.Listen(tcp, domain)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// accept conn
		conn, err := lis.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Create a new client
		client := &Client{
			conn:   conn,
			writer: bufio.NewWriter(conn),
		}

		// Register the new client
		clients[client] = true

		go handleConn(client)

		go broadcaster()

	}

}

func handleConn(client *Client) {
	// defer func() {
	// 	// Unregister the client when the connection is closed
	// 	leave <- client
	// 	client.conn.Close()
	// }()
	frame := lib.NewFrame()

	// read first http message
	buf := make([]byte, 1024)
	n, err := client.conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf[:n]))

	// get transformed httpy
	lines := string(buf[:n])
	httpy := lib.NewHTTPy().GetHTTPy(lines)
	lib.PrintJSON(httpy)
	key := httpy.SecWebKey[1:]

	// get upgrade response string
	responseHTTPString := lib.GetUpgradeResponseString(key)
	fmt.Println("response: ")
	fmt.Println(string(responseHTTPString))

	// write http response to client
	_, err = client.conn.Write(responseHTTPString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("websocket upgrade success!\n\n")

	// get first message
	buf2 := make([]byte, 4096)
	_, err = client.conn.Read(buf2)
	if err != nil {
		fmt.Println("read buffer failed")
		log.Fatal(err)
	}
	message := frame.DecodeFrame(buf2)
	fmt.Printf("[*]client: ")
	fmt.Println(string(message))

	// response first message
	rawMsg := []byte("welcome to qchat")
	encodedMsg := frame.EncodeFrame(rawMsg)

	// _, err = client.conn.Write(encodedMsg)
	// if err != nil {
	// 	fmt.Println("write buffer failed")
	// 	log.Fatal(err)
	// }

	// broadcast message
	broadcast <- encodedMsg

}

func broadcaster() {
	for {
		select {
		case message := <-broadcast:
			for client := range clients {
				fmt.Println("casting ")
				// Send the message to each connected client
				_, err := client.writer.Write(message)
				if err != nil {
					log.Printf("Error broadcasting message: %v", err)
				}
				client.writer.Flush()
			}
		}
	}
}
