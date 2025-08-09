package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func getUDPAddr() *net.UDPAddr {

	address, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return address
}

func getUDPConn() *net.UDPConn {
	conn, err := net.DialUDP("udp", nil, getUDPAddr())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return conn
}

func main() {

	conn := getUDPConn()
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Print(err)
		}
		conn.Write([]byte(text))
	}

}
