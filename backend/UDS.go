package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

const socketPath = "/tmp/clipboard.sock"

func connection(){
    log.Println("Establishing connection with USD...")

	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
    if err != nil {
        panic(err)
    }
    log.Println("Connection established.")
	defer listener.Close()

	for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Accept error:", err)
            continue
        }
        log.Println("Frontend connected.")
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
    switch received {
        case "get_history":
            log.Println("Request received from frontend.")
            response,err := json.Marshal(getClipboardHistory())
            if err == nil{
                conn.Write([]byte(response))
                log.Println("Response sent.")
            }
    }
}