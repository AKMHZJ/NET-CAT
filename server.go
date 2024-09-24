package main

import (
	"fmt"
	"net"
	"os"
	"sync"

	con "netcat/functions"
)

const defaultport = 8989

func main() {
	port := defaultport
	if len(os.Args) == 2 {
		costumport, err := con.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		port = costumport
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	fmt.Printf("Starting Server at port :%d\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error Starting a server : ", err)
	}
	defer listener.Close()

	Clients := make(map[*con.Client]bool)
	var ClientMutex sync.Mutex
	Maxclients := 10

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection : ", err)
			continue
		}

		ClientMutex.Lock()
		if len(Clients) >= Maxclients {
			conn.Write([]byte("Sorry, the chat room is full. Please try again later.\n"))
			conn.Close()
			ClientMutex.Unlock()
			continue
		}

		client := con.Client{Conn: conn}
		Clients[&client] = true
		ClientMutex.Unlock()

		go con.Hundleclients(&client, Clients)
	}
}
