package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/Waterbootdev/http-from-tcp/internal/request"
	"github.com/Waterbootdev/http-from-tcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{listener: listener, handler: handler}

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

	request, err := request.RequestFromReader(conn)

	log.Printf("Request: %v", request)

	if err != nil {
		headers := response.GetDefaultHeaders(0)
		response.WriteStatusLine(conn, response.BAD_REQUEST)
		response.WriteHeaders(conn, headers)
		log.Printf("Connection from %s closed", conn.RemoteAddr())
		return
	}

	buffer := &bytes.Buffer{}

	handlerErr := s.handler(buffer, request)

	if handlerErr != nil {
		handlerErr.Write(conn)
		return
	}

	headers := response.GetDefaultHeaders(buffer.Len())
	response.WriteStatusLine(conn, response.OK)
	response.WriteHeaders(conn, headers)

	conn.Write(buffer.Bytes())

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
