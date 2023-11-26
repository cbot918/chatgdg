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
			defer conn.Close()
			// read first http message
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(buf[:n]))

			// get httpy
			lines := string(buf[:n])
			httpy := NewHTTPy().GetHTTPy(lines)

			PrintJSON(httpy)
			fmt.Println("here")
			fmt.Println(httpy.SecWebKey[1:])

			// get upgrade response string
			upgrader := NewWsUpgrader()
			temp := httpy.SecWebKey[1:]
			fmt.Println("key:", temp)
			fmt.Println("key length: ", len(temp))
			key := upgrader.getReturnSec(temp)
			fmt.Println("output key:", key)
			fmt.Println("put[ut key length: ", len(key))
			resString := upgrader.GetResponseString(key)

			_, err = conn.Write(resString)
			if err != nil {
				log.Fatal(err)
			}
		}(conn)

	}

}
