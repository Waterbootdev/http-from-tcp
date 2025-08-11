package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Waterbootdev/http-from-tcp/internal/request"
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

		request, err := request.RequestFromReader(conn)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(request.RequestLine.RequestLineString())
		fmt.Println()
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
