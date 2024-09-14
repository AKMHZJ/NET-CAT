package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

var MessageHistory []string

// HistoryMutex   sync.Mutex

type Client struct {
	name string
	conn net.Conn
}

const defaultport = 8989

func Hundleclients(client *Client, Clients map[*Client]bool) {
	defer client.conn.Close()

	client.conn.Write([]byte("Welcome to TCP-Chat!\n"))

	data, _ := os.ReadFile("pinguen.txt")
	// if err != nil {
	// 	fmt.Println("error reading the pinguen.")
	// 	return
	// }
	data = append(data, byte('\n'))
	_, err := client.conn.Write(data)

	client.conn.Write([]byte("[ENTER YOUR NAME]: "))
	name := make([]byte, 64)
	n, err := client.conn.Read(name)
	if err != nil {
		fmt.Print("Error reading client name: ", err)
	}

	client.name = string(name[:n-1])

	// HistoryMutex.Lock()
	for _, msg := range MessageHistory {
		client.conn.Write([]byte(msg))
	}
	// HistoryMutex.Unlock()

	message := fmt.Sprintf("%s has joined our chat...\n", client.name)
	BroadcastingMsg(message, client, Clients)

	for {
		Message := make([]byte, 1024)

		z, err := client.conn.Read(Message)
		if err != nil {
			if err == io.EOF {
				delete(Clients, client)
				message := fmt.Sprintf("%s has left our chat...\n", client.name)
				BroadcastingMsg(message, client, Clients)
				continue
			}
		}
		msg := string(Message[:z-1])
		for _, m := range msg {
			if m < 32 {
				msg = ""
				break
			}
		}

		formattedmsg := fmt.Sprintf("[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), client.name, msg)
		BroadcastingMsg(formattedmsg, client, Clients)
	}
}

func BroadcastingMsg(message string, sender *Client, Clients map[*Client]bool) {
	// HistoryMutex.Lock()
	// defer HistoryMutex.Unlock()

	// HistoryMutex.Lock()
	// MessageHistory = append(MessageHistory, message)
	// HistoryMutex.Unlock()

	for client := range Clients {
		if client != sender {
			if client.name != "" {
				_, err := client.conn.Write([]byte(message))
				if err != nil {
					fmt.Println("Error broadcasting message : ", err)
				}
			}
		}
	}
}

func main() {
	port := defaultport
	if len(os.Args) > 1 {
		costumport, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
		port = costumport
	}

	fmt.Printf("Starting Server at port :%d\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error Starting a server : ", err)
	}
	defer listener.Close()

	Clients := make(map[*Client]bool)
	// var ClientMutex sync.Mutex
	Maxclients := 10

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection : ", err)
			continue
		}

		if len(Clients) >= Maxclients {
			conn.Write([]byte("Sorry, the chat room is full. Please try again later.\n"))
			conn.Close()
			continue
		}

		client := &Client{conn: conn}
		Clients[client] = true

		go Hundleclients(client, Clients)
	}
}
