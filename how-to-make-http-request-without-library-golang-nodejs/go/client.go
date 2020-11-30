package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func ReadAllBuffer(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 32000)

	tot := 0
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err == io.EOF {
				fmt.Println(err)
				return buf, nil
			}
			return nil, err
		}

		buf = append(buf, tmp[:n]...)
		fmt.Println(n)
		tot = tot + n
		if tot == 1000000 {
			return buf, nil
		}

		time.Sleep(2*time.Second)
	}
}


func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()

	if _, err := conn.Write([]byte("hello")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("about to read")

	n, err := ReadAllBuffer(conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(len(n))
	//fmt.Println(string(buff))
}
