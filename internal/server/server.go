package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

const HELLO_WORLD_RESPONSE = `HTTP/1.1 200 OK
Content-Type: text/plain

Hello World!`

const HELLO_WORLD = "Hello World!"

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{listener: listener}

	go func() {
		server.listen()
	}()

	return server, nil
}

func (s *Server) Close() error {
	if s.closed.Swap(true) {
		return nil
	}

	log.Printf("server closed")

	return s.listener.Close()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	log.Printf("Handling connection from %s", conn.RemoteAddr())
	headers := response.GetDefaultHeaders(len(HELLO_WORLD))
	response.WriteStatusLine(conn, response.OK)
	response.WriteHeaders(conn, headers)
	conn.Write([]byte(HELLO_WORLD))
	log.Printf("Connection from %s closed", conn.RemoteAddr())
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
		go s.handle(conn)
	}
}
