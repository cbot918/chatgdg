type Client struct {
	conn   net.Conn
	writer *bufio.Writer
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan []byte)
	// join      = make(chan *Client)
	// leave     = make(chan *Client)
)

func main() {

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
			/* *********************** */

			/* Add client to clients */
			client := &Client{
				conn:   conn,
				writer: bufio.NewWriter(conn),
			}
			// defer func() {
			// 	// Unregister the client when the connection is closed
			// 	leave <- client
			// 	client.conn.Close()
			// }()

			clients[client] = true
			// join <- client
			fmt.Println(client)
			fmt.Printf("joined\n")
			fmt.Printf("users: %d\n\n", len(clients))
			/* ***********************8 */

			for {
				// /* WS Read Write */
				data, err = read(conn)
				if err != nil {
					continue
				}
				fmt.Printf("received : %q\n", string(decodeFrame(data)))

				broadcast <- data

			}

		}(conn)

		go handleMessages()

	}

}

func handleMessages() {
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
