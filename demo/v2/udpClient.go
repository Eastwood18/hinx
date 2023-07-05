package main

import (
	"fmt"
	"net"
	"os"
	//  "io"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:8899")
	defer conn.Close()
	if err != nil {
		os.Exit(1)
	}
	fmt.Println([]byte("Hello world!"))

	conn.Write([]byte("Hello world!"))

	fmt.Println("send msg")

	//var msg [20]byte
	//conn.Read(msg[0:])
	//
	//fmt.Println("msg is", string(msg[0:10]))
}
