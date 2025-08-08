package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("connection has been accepted")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
	}
}
