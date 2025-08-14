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

	if err != nil {
		(&HandlerError{StatusCode: response.BAD_REQUEST, Message: err.Error()}).Write(conn)
		return
	}

	buffer := &bytes.Buffer{}

	handlerErr := s.handler(buffer, request)

	if handlerErr != nil {

		handlerErr.Write(conn)

		return
	}

	response.WriteBufferOk(conn, buffer)

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
