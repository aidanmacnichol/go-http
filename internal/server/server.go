package server

import (
	"fmt"
	"go-http/internal/request"
	"go-http/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	Closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		Listener: listener,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.Closed.Store(true)
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.Closed.Load() {
				return
			}
			log.Printf("Accept Error: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	_, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("Error parsing request: %v\n", err.Error())
		return
	}
	response.WriteStatusLine(conn, response.StatusOK)
	headers := response.GetDefaultHeaders(0)
	if err := response.WriteHeaders(conn, headers); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}
