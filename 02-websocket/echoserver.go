package main

import (
	"fmt"
	"os"
	"net"
)

func main() {
	fmt.Println("hello responder")
	fmt.Println("you supplied", len(os.Args), "arguments")
	if len(os.Args) != 2 {
		fmt.Println("Usage: responder <port>")
		os.Exit(0)
	}
	port := os.Args[1]
	fmt.Println("Port:", port)
	
	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("ip4", service)
	checkError(err)
	fmt.Println("Resolved TCP address for the service")
	
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		
		go ServeSocket(conn)
	}
}

func ServeSocket(conn *net.TCPConn) {
	fmt.Println("In ServeSocket")
	for {
		var buffer [512]byte
		n, err := conn.Read(buffer[:])
		if err != nil {
			break
		}
		message := string(buffer[:n])
		fmt.Println("Got message:", message)
		if message[:4] == "quit" {
			break
		}
		_, err = conn.Write([]byte(message))
		if err != nil {
			break
		}
	}
	conn.Close()
	fmt.Println("Connection closed")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

