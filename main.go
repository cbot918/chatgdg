package main

import (
	"chatdgd/src/util"
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	protocol = "tcp"
	host     = "localhost:8888"
)

var mutex sync.Mutex

func main() {

	fmt.Println("listening: ", host)
	listener, err := net.Listen(protocol, host)
	if err != nil {
		log.Fatal(err)
	}

	ws := NewWebsocket()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go ws.handleWSConn(conn)

	}

}

type Client struct {
	Name string
	Conn net.Conn
}

func (c *Client) read() (string, error) {
	buf := make([]byte, 1024)
	n, err := c.Conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

type Websocket struct {
	Count int
	Conns map[int]*Client
}

func NewWebsocket() *Websocket {
	return &Websocket{
		Conns: make(map[int]*Client),
	}
}

const (
	GET  = "GET"
	POST = "POST"
)

func (w *Websocket) handleWSConn(conn net.Conn) {

	w.Count += 1

	mutex.Lock()
	w.Conns[w.Count] = &Client{
		Name: "yale",
		Conn: conn,
	}
	mutex.Unlock()

	fmt.Printf("number of current users: %d\n", w.Count)
	fmt.Printf("%v connected\n", conn)

	// fmt.Println(conn)
	// fmt.Printf("%#+v", conn)
	// printJSON(conn.RemoteAddr())

	lines, err := w.Conns[1].read()
	if err != nil {
		fmt.Println(err)
		// log.Fatal("read line failed")
	}

	httpy := util.GetHTTPy(lines)

	fmt.Println("source http content: ")
	fmt.Println(lines)

	fmt.Println("httpy data: ")
	util.PrintJSON(httpy)

	upgrader := util.NewWsUpgrader()

	err = upgrader.Upgrade(w.Conns[1].Conn, httpy.SecWebKey)
	if err != nil {
		fmt.Println(err)
	}
}
