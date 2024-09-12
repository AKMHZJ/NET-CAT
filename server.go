package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	messageHistory []string
	historyMutex   sync.Mutex
)

const defaultPort = 8989

type Client struct {
	name string
	conn net.Conn
}

func Handlclient(client *Client, Clients map[*Client]bool, clientsMutex *sync.Mutex) {
	defer client.conn.Close()
	
	client.conn.Write([]byte("Welcome to TCP-Chat!\n"))

	data, _ := os.ReadFile("ascii.txt")
	data = append(data, byte('\n'))
	_, err := client.conn.Write(data)

	client.conn.Write([]byte("ENTER YOUR NAME : "))
	name := make([]byte, 64)
	n, err := client.conn.Read(name)
	if err != nil {
		fmt.Println("Error reading client name: ", err)
		return
	}
	client.name = string(name[:n-1])

	historyMutex.Lock()
	for _, msg := range messageHistory {
		client.conn.Write([]byte(msg))
	}
	historyMutex.Unlock()

	// announce the new client
	Message := fmt.Sprintf("[%s] #%s# has joined our chat...\n", time.Now().Format("2006-01-02 15:04:05"), client.name)
	Broadcastmessage(Message, Clients, client, clientsMutex)

	for {
		message := make([]byte, 1024)
		n, err := client.conn.Read(message)
		if err != nil {
			if err == io.EOF {
				delete(Clients, client)
				message := fmt.Sprintf("[%s] #%s# has left our chat...\n", time.Now().Format("2006-01-02 15:04:05"), client.name)
				Broadcastmessage(message, Clients, client, clientsMutex)
				return
			}
			fmt.Println("Error reading from client : ", err)
			continue
		}
		msg := string(message[:n-1])
		if msg == "" {
			continue
		}

		Formattedmsg := fmt.Sprintf("[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), client.name, msg)
		Broadcastmessage(Formattedmsg, Clients, client, clientsMutex)
	}
}

func Broadcastmessage(Message string, Clients map[*Client]bool, sender *Client, clientsMutex *sync.Mutex) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	historyMutex.Lock()
	messageHistory = append(messageHistory, Message)
	defer historyMutex.Unlock()

	for client := range Clients {
		if client != sender {
			_, err := client.conn.Write([]byte(Message))
			if err != nil {
				fmt.Println("Error broadcasting massage", err)
			}
		}
	}
}

func main() {
	port := defaultPort
	if len(os.Args) > 1 {
		customPort, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
		port = customPort
	}

	fmt.Printf("Server is listening on port %d\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	Clients := make(map[*Client]bool)
	var clientsMutex sync.Mutex
	Maxclient := 10

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		clientsMutex.Lock()
		if len(Clients) >= Maxclient {
			clientsMutex.Unlock()
			conn.Write([]byte("Sorry, the chat room is full. Please try again later.\n"))
			conn.Close()
			continue
		}

		client := &Client{conn: conn}
		Clients[client] = true
		clientsMutex.Unlock()

		go Handlclient(client, Clients, &clientsMutex)
	}

	// http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

// 1: Set up the basic server structure
// 2: TODO: Implement connection handling
