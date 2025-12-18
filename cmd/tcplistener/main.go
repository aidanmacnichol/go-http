package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const portAddr = ":42069"

func getLinesChannel(conn net.Conn) <-chan string {
	lineCh := make(chan string)
	go func() {
		defer conn.Close()
		defer close(lineCh)
		lineStr := ""
		for {
			buf := make([]byte, 8, 8)
			n, err := conn.Read(buf)
			if err != nil {
				if lineStr != "" {
					lineCh <- lineStr
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}

			str := string(buf[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lineCh <- fmt.Sprintf("%s%s\n", lineStr, parts[i])
				lineStr = ""
			}
			lineStr += parts[len(parts)-1]
		}
	}()
	return lineCh
}

func main() {
	listener, err := net.Listen("tcp", portAddr)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", portAddr, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Could not open %s: %s\n", portAddr, err)
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())
		lineCh := getLinesChannel(conn)
		for line := range lineCh {
			fmt.Printf(" %s", line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
