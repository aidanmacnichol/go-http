package main

import (
	"fmt"
	"log"
	"net"

	"go-http/internal/request"
)

const portAddr = ":42069"

func main() {
	listener, err := net.Listen("tcp", portAddr)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", portAddr, err)
	}
	defer listener.Close()

	fmt.Println("Listening for traffic on", portAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Could not open %s: %s\n", portAddr, err)
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error parsing request: %v", err.Error())
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)
	}
}
