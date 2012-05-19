package main

import (
	"fmt"
	"net/http"
	"os"
	"code.google.com/p/go.net/websocket"
	"bufio"
	"net"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: wstelnet port")
		os.Exit(0)
	}
	
	port := os.Args[1]
	fmt.Println("Serving web on port", port)
	service := ":" + port
	
	http.Handle("/script/", http.FileServer(http.Dir(".")))
	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/", http.FileServer(http.Dir("./html/")))
	http.Handle("/websocket/", websocket.Handler(ProcessSocket))
	err := http.ListenAndServe(service, nil)
	checkError(err)
}

func ProcessSocket(ws *websocket.Conn) {
	fmt.Println("In ProcessSocket")
	var msg string

	err := websocket.Message.Receive(ws, &msg)
	if err != nil {
		fmt.Println("ProcessSocket: got error", err)
		_ = websocket.Message.Send(ws, "FAIL:" + err.Error())
		return
	}
	fmt.Println("ProcessSocket: got message", msg)

	service := msg
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Println("Error in ResolveTCPAddr:", err)
		_ = websocket.Message.Send(ws, "FAIL:" + err.Error())
		return
	}
	
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Error in DialTCP:", err)
		_ = websocket.Message.Send(ws, "FAIL:" + err.Error())
		return
	}
	
	_ = websocket.Message.Send(ws, "SUCC")

	RunTelnet(ws, conn)
}

func RunTelnet(ws *websocket.Conn, conn *net.TCPConn) {
	fmt.Println("Running telnet")
	go ReadSocket(ws, conn)
	
	// Read websocket and write to socket.
	crlf := []byte{13, 10}
	var msg string
	for {
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			_ = conn.Close()
			break
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			break
		}
		fmt.Println("Sent message to host:", msg)
		// Send \r\n (as HTTP protocol requires)
		_, err = conn.Write(crlf)
		if err != nil {
			break
		}
	}
	fmt.Println("RunTelnet exit")
}

// Read from socket and write to websocket
func ReadSocket(ws *websocket.Conn, conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	var line string = ""
	for {
		if reader == nil {
			break
		}
		buffer, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		// fmt.Println("ReadSocket: got", len(buffer), "bytes")
		line = line + string(buffer)
		if !isPrefix {
			// fmt.Println("Sending message to web socket:", line)
			err = websocket.Message.Send(ws, line)
			if err != nil {
				_ = conn.Close()
			}
			line = ""
		}
	}
	fmt.Println("ReadSocket exit")
	ws.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

