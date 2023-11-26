package main

import (
	"fmt"
	"log"
	"net"
)

const (
	tcp    = "tcp"
	domain = ":8888"
)

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

		go func(conn net.Conn) {

			// read first http message
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(buf[:n]))

			// get transformed httpy
			lines := string(buf[:n])
			httpy := NewHTTPy().GetHTTPy(lines)
			PrintJSON(httpy)
			key := httpy.SecWebKey[1:]

			// get upgrade response string
			responseHTTPString := GetUpgradeResponseString(key)
			fmt.Println("response: ")
			fmt.Println(string(responseHTTPString))

			// write http response to client
			_, err = conn.Write(responseHTTPString)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("websocket upgrade success!")

			// get first message
			buf2 := make([]byte, 4096)
			_, err = conn.Read(buf2)
			if err != nil {
				fmt.Println("read buffer failed")
				log.Fatal(err)
			}

			f := NewFrame()
			message := f.DecodeFrame(buf2)
			fmt.Printf("%q", string(message))
		}(conn)

	}

}
