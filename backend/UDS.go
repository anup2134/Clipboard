package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const socketPath = "/tmp/clipboard.sock"

func connection(){
	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
    if err != nil {
        panic(err)
    }
	defer listener.Close()

	for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Accept error:", err)
            continue
        }

        go handleConnection(conn)
    }
}


func handleConnection(conn net.Conn) {
    defer conn.Close()

    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Read error:", err)
        return
    }

    received := string(buf[:n])

	if received == "get_history"{
		response,err := json.Marshal(getClipboardHistory())
		if err == nil{
			conn.Write([]byte(response))
		}
	}
}