package main

import (
	"bufio"
	"log"
	"net"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var openConnection = make(map[net.Conn]bool)
var newConnection = make(chan net.Conn)
var deadConnection = make(chan net.Conn)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	logFatal(err)
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			logFatal(err)

			openConnection[conn] = true
			newConnection <- conn
		}
	}()
	for {
		select {
		case conn := <-newConnection:
			go broadcastMessage(conn)
		case conn := <-deadConnection:
			for item := range openConnection {
				if item == conn {
					break
				}
			}
			delete(openConnection, conn)
		}
	}
}
func broadcastMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')

		if err != nil {
			break
		}
		for item := range openConnection {
			if item != conn {
				item.Write([]byte(message))
			}
		}
	}
	deadConnection <- conn
}
