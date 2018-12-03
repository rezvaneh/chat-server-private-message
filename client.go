package main

import (
	"bufio"
	"chat-server-private-message/config"
	"fmt"
	"net"
	"os"
	"strings"
)

type Command struct {
	Command, Username, Body string
}

func main() {
	username, configs := getConfig()
	conn, err := net.Dial("tcp", configs.Hostname+":"+configs.Port)
	config.CheckForError(err, "Connection refused")
	defer conn.Close()

	go handleConnectionInput(username, conn)
	getConsoleInput(conn)
}

func getConfig() (string, config.Configs) {
	if len(os.Args) >= 2 {
		username := os.Args[1]
		configs := config.LoadConfig()
		return username, configs
	} else {
		println("Please provide a username as the first parameter.")
		os.Exit(1)
		return "", config.Configs{}
	}
}

func handleConnectionInput(username string, conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		config.CheckForError(err, "Lost server connection")
		message = strings.TrimSpace(message)
		if message != "" {
			Command := getServerCommandParts(message)
			switch Command.Command {
			case "start":
				sendCommand("new", username, conn)

			case "connect":
				fmt.Printf("Welcome %s!"+"\n", username)
				fmt.Printf("Please Enter Your Friend's Usernames:" + "\n")

			case "getUsername":
				fmt.Printf("Please Enter Your Friend's Usernames:" + "\n")

			case "getMessage":
				fmt.Printf("Please Enter Your Message: \n")

			case "sendMessage":
				fmt.Printf("[%s]: %s \n", Command.Username, Command.Body)

			case "error":
				fmt.Printf("Duplicated username. Please try another one. \n")
			}
		}
	}
}

func getConsoleInput(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		config.CheckForError(err, "Lost console connection")

		message = strings.TrimSpace(message)
		if message != "" {
			command := getServerInputParts(message)
			switch command.Command {
			case "":
				sendCommand("message", message, conn)

			case "getUsername":
				sendCommand("getMessage", command.Username, conn)
			}
		}
	}
}

func sendCommand(command string, body string, conn net.Conn) {
	message := fmt.Sprintf("/%v %v\n", command, body)
	conn.Write([]byte(message))
}

func getServerCommandParts(message string) Command {
	s := strings.Split(message, " ")
	result := strings.Join(s[2:], " ")
	if len(s) > 0 {
		return Command{
			Command:  strings.Replace(s[0], "/", "", 1),
			Username: s[1],
			Body:     result,
		}
	}
	return Command{}
}

func getServerInputParts(message string) Command {
	s := strings.Split(message, " ")
	if len(s) > 1 {
		result := strings.Join(s, " ")
		if string(message[0]) != "/" {
			return Command{
				Body: result,
			}
		}

		return Command{
			Command: s[0],
			Body:    result,
		}
	}

	return Command{
		Body: message,
	}
}
