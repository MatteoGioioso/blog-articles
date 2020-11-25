package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		
		
		fmt.Println("connection received")
		if _, err := io.WriteString(conn, makeBigMessage()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func makeBigMessage() string {
	return string(make([]byte, 1000000))
}
