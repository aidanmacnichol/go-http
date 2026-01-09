package server

import (
	"fmt"
	"go-http/internal/request"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
)

// Contains the state of the server
type Server struct {
	Listener    net.Listener
	Closed      uint32 // Bool flag - 32 bit because many CPU architectures only support world aligned boundaries
	mu          sync.Mutex
	ActiveConns map[net.Conn]bool
}

// Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
func Serve(port int) (*Server, error) {
	portAddr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", portAddr)
	if err != nil {
		return nil, err
	}
	s := Server{
		Listener: listener,
		Closed:   0,
	}
	go s.listen()
	return &s, nil
}

// Closes the listener and the server
func (s *Server) Close() error {
	atomic.StoreUint32(&s.Closed, 1)
	fmt.Println("Closing server...")
	s.Listener.Close()

	// What to do with errors????
	return nil
}

// Uses a loop to .Accept new connections as they come in,
// and handles each one in a new goroutine. I used an atomic.
// Bool to track whether the server is closed or not so that I can ignore connection errors after the server is closed.
func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Accept Error: ", err)
			continue
		}

		s.mu.Lock()
		s.ActiveConns[conn] = true
		s.mu.Unlock()

		go s.handle(conn)
	}
}

// Handles a single connection by writing the following response and then closing the connection:
func (s *Server) handle(conn net.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.ActiveConns, conn)
		s.mu.Unlock()
		conn.Close()
	}()

	_, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("Error parsing request: %v\n", err.Error())
		return
	}

	fmt.Println("HTTP/1.1 OK")
	fmt.Println("Content-Type: text/plain")
	fmt.Println("\nHello World!")
}
