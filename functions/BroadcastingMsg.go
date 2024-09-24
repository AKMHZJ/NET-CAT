package netcat

import (
	"fmt"
	"time"
)

func BroadcastingMsg(message string, sender *Client, Clients map[*Client]bool) {
	HistoryMutex.Lock()
	defer HistoryMutex.Unlock()

	MessageHistory = append(MessageHistory, message)

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
