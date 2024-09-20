package netcat

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	MessageHistory []string
	HistoryMutex   sync.Mutex
)

type Client struct {
	Name string
	Conn net.Conn
}

func Hundleclients(client *Client, Clients map[*Client]bool) {
	defer client.Conn.Close()

	client.Conn.Write([]byte("Welcome to TCP-Chat!\n"))

	data, err := os.ReadFile("penguin.txt")
	if err != nil {
		fmt.Println("error reading the pinguen.")
		return
	}
	data = append(data, byte('\n'))
	client.Conn.Write(data)
start:
	client.Conn.Write([]byte("[ENTER YOUR NAME]: "))
	name := make([]byte, 64)
	n, err := client.Conn.Read(name)
	if err != nil {
		fmt.Print("Error reading client name: ", err)
	}

	NAME := string(name[:n-1])
	Skiip := false
	for _, n := range NAME {
		if n < 32 {
			NAME = ""
			Skiip = true
		}
	}
	if Skiip || strings.ReplaceAll(NAME, " ", "") == "" {
		client.Conn.Write([]byte("\033[31;1;4minvalid message !!!\033[0m\n"))
		goto start
	}

	for names := range Clients {
		if string(name[:n-1]) == names.Name {
			client.Conn.Write([]byte("\033[31;1;4mthis name is already exist !!!\033[0m\n"))
			goto start
		}
	}
	client.Name = string(name[:n-1])

	HistoryMutex.Lock()
	for _, msg := range MessageHistory {
		client.Conn.Write([]byte(msg))
	}
	HistoryMutex.Unlock()
	msg := fmt.Sprintf("\n%s has joined our chat...\n", client.Name)
	BroadcastingMsg(msg, client, Clients)
	Message := make([]byte, 1024)
	for {
		if Message[1023] == 0 || Message[1023] == 12 {
			formattedmInput := fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), client.Name)
			client.Conn.Write([]byte(formattedmInput))
		}
		Message = make([]byte, 1024)
		z, err := client.Conn.Read(Message)
		if err != nil {
			if err == io.EOF {
				delete(Clients, client)
				message := fmt.Sprintf("\n%s has left our chat...\n", client.Name)
				BroadcastingMsg(message, client, Clients)
				return
			}
		}
		msg := string(Message[:z-1])
		Skiip := false
		for _, m := range msg {
			if m < 32 {
				msg = ""
				Skiip = true
			}
		}
		if Skiip || strings.ReplaceAll(msg, " ", "") == "" {
			client.Conn.Write([]byte("\033[31;1;4minvalid message !!!\033[0m\n"))
			continue
		}
		formattedmsg := fmt.Sprintf("\n[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), client.Name, msg)
		BroadcastingMsg(formattedmsg, client, Clients)
	}
}

func BroadcastingMsg(message string, sender *Client, Clients map[*Client]bool) {
	HistoryMutex.Lock()
	defer HistoryMutex.Unlock()

	// HistoryMutex.Lock()
	MessageHistory = append(MessageHistory, message)
	// HistoryMutex.Unlock()

	for client := range Clients {
		if client != sender {
			if client.Name != "" {
				_, err := client.Conn.Write([]byte(message))
				formattedInput := fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), client.Name)
				client.Conn.Write([]byte(formattedInput))
				if err != nil {
					fmt.Println("Error broadcasting message : ", err)
					continue
				}
			}
		}
	}
}
