package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Client chan<- string

var (
	incommingClients = make(chan Client)
	leavingClients   = make(chan Client)
	messages         = make(chan string)
)

var (
	host = flag.String("h", "localhost", "host")
	port = flag.Int("p", 3090, "port")
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	message := make(chan string)
	go MessageWrite(conn, message)

	clientName := conn.RemoteAddr().String()
	message <- fmt.Sprintf("Welcome to the server %s\n", clientName)
	messages <- fmt.Sprintf("New client is here, say welcome to %s\n", clientName)
	incommingClients <- message

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- fmt.Sprintf("%s: %s\n", clientName, input.Text())
	}

	leavingClients <- message
	messages <- fmt.Sprintf("%s said goodbye!", clientName)
}

func MessageWrite(conn net.Conn, messages <-chan string) {
	for message := range messages {
		fmt.Fprintln(conn, message)
	}
}
