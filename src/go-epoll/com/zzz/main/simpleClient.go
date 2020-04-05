package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	defer conn.Close()
	if err != nil {
		fmt.Println("net.Dial err=", err)
	}
	conn.Write([]byte("hello world"))
	time.Sleep(time.Second * 2)
}
