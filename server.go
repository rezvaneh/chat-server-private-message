package main

import (
	"bufio"
	"chat-server-private-message/config"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	Connection net.Conn
	Username   string
	Friend     string
	Configs    config.Configs
}

type Message struct {
	SenderUsername   string
	ReceiverUsername string
	Text             string
}

var clients []*Client

func main() {
	configs := config.LoadConfig()
	psock, err := net.Listen("tcp", ":"+configs.Port)
	config.CheckForError(err, "Server can't start.")
	fmt.Printf("start... port %v \n", configs.Port)

	for {
		conn, _ := psock.Accept()
		client := Client{Connection: conn, Configs: configs}
		client.Register()

		var msg Message
		channel := make(chan string)

		go getInput(channel, &client)
		go sendResponse(channel, &client, &msg)
		SendMessageClient("start", configs.Port, &client, "")
	}
}

func getInput(send chan string, client *Client) {
	defer close(send)

	reader := bufio.NewReader(client.Connection)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("err", err)
			client.Close()
			return
		}
		send <- string(line)
	}
}

func sendResponse(receive <-chan string, client *Client, msg *Message) {
	for {
		message := <-receive
		message = strings.TrimSpace(message)
		cmd, body := getCommandParts(message)
		if cmd != "" {
			switch cmd {
			case "new":
				err := client.checkUsername(body)
				if err != nil {
					SendMessageClient("error", body, client, "")
					client.Close()
					return
				}
				client.Username = body
				msg.SenderUsername = client.Username
				SendMessageClient("connect", body, client, "")

			case "message":
				if msg.SenderUsername == "" {
					SendMessageClient("getUsername", body, client, "")

				} else if msg.ReceiverUsername == "" {
					msg.ReceiverUsername = body
					SendMessageClient("getMessage", body, client, "")

				} else if msg.Text == "" {
					msg.Text = body
					SendMessageClient("sendMessage", body, client, msg.ReceiverUsername)
					SendMessageClient("getUsername", body, client, "")
					msg.ReceiverUsername = ""
					msg.Text = ""
				} else {
					msg.ReceiverUsername = ""
					msg.Text = ""
					SendMessageClient("getUsername", body, client, "")
				}
			}
		}
	}
}

func getCommandParts(message string) (string, string) {
	s := strings.Split(message, " ")
	if len(s) > 0 {
		result := strings.Join(s[1:], " ")
		return strings.Replace(s[0], "/", "", 1), result
	}
	return "", ""
}

func SendMessageClient(messageType string, message string, client *Client, friend string) {
	for _, _client := range clients {
		if _client.Username == friend && friend != "" {
			payload := fmt.Sprintf("/%v %v %v", messageType, client.Username, message)
			fmt.Fprintln(_client.Connection, payload)
			break
		} else if _client.Username == client.Username && messageType != "sendMessage" {
			payload := fmt.Sprintf("/%v %v %v", messageType, client.Username, message)
			fmt.Fprintln(_client.Connection, payload)
			break
		}
	}
}

func (client *Client) Register() {
	clients = append(clients, client)
}

func (client *Client) Close() {
	client.Connection.Close()
	for i := 0; i < len(clients); i++ {
		if client == clients[i] {
			clients = append(clients[:i], clients[i+1:]...)
		}
	}
}

func (client *Client) checkUsername(username string) error {
	for _, c := range clients {
		if c.Username == username {
			return errors.New("error")
		}
	}
	return nil
}
