package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

func read(fd int, res *[]byte)  {
	buff := make([]byte, 10)

	for {
		tmp := make([]byte, 10)
		n, err := unix.Read(fd, tmp)
		if n == 0 && err == nil {
			// EOF
			break
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		buff = append(buff, tmp...)
	}
	*res = buff
}

func main() {
	// domain: protocol family to be used in the communication. AF_INET: IPv4 family
	// type: SOCK_STREAM is TCP socket
	// proto: (Protocol) refers to a single protocol that supports the selected socket.
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println(err)
	}

	sa := &unix.SockaddrInet4{Port: 9000}
	if err := unix.Connect(fd, sa); err != nil {
		fmt.Println(err)
	}

	req := fmt.Sprint("GET / HTTP/1.1\nHost: localhost:9000\nConnection: close\n\n")

	if _, err := unix.Write(fd, []byte(req)); err != nil {
		fmt.Println(err)
	}

	res := make([]byte, 1024)

	read(fd, &res)

	fmt.Println(string(res))

	if err := unix.Close(fd); err != nil {
		fmt.Println(err)
	}
}
