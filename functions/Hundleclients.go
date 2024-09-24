package netcat

import (
	"fmt"
	"io"
	"log"
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

	log.Printf("New client connection from %s", client.Conn.RemoteAddr())

	client.Conn.Write([]byte("Welcome to TCP-Chat!\n"))

	data, err := os.ReadFile("penguin.txt")
	if err != nil {
		log.Printf("Error reading penguin.txt: %v", err)
		return
	}
	data = append(data, byte('\n'))
	client.Conn.Write(data)
start:
	client.Conn.Write([]byte("[ENTER YOUR NAME]: "))
	name := make([]byte, 64)
	n, err := client.Conn.Read(name)
	if err != nil {
		log.Printf("Error reading client name: %v", err)
		fmt.Println("Error reading client name: ", err)
	}
	if n < 1 {
		return
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

	log.Printf("Client %s has joined the chat", client.Name)

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
				log.Printf("Client %s has left the chat", client.Name)
				return
			}
			log.Printf("Error reading message from %s: %v", client.Name, err)
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
		} else if strings.HasPrefix(msg, "/username") {
			oldName := client.Name
			n := strings.Join(strings.Fields(msg)[1:], " ")
			tf := true
			for names := range Clients {
				if n == names.Name {
					client.Conn.Write([]byte("\033[31;1;4mthis name is already exist !!!\033[0m\n"))
					tf = false
				}
			}
			if tf && n != "" {
				client.Name = n
				client.Conn.Write([]byte("you change your name to " + client.Name + " ...\n"))
				message := fmt.Sprintf("\n%s has change his name to %s ...\n", oldName, client.Name)
				BroadcastingMsg(message, client, Clients)
				log.Printf("Client %s changed their name to %s", oldName, client.Name)
			} else {
				client.Conn.Write([]byte("\033[31;1;4mplease write your name after the flag :\033[0m\n"))
				client.Conn.Write([]byte("\033[31;1;4m[USAGE]: /username [$the new name] .\033[0m\n"))
			}
			continue
		}
		formattedmsg := fmt.Sprintf("\n[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), client.Name, msg)
		BroadcastingMsg(formattedmsg, client, Clients)
		log.Printf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), client.Name, msg)
	}
}
