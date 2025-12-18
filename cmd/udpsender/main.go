package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const portAddr = ":42069"

func main() {
	addr, err := net.ResolveUDPAddr("udp", portAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving address: %v", err)
		os.Exit(1)
	}

	udpConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Error dailing address: %v", err)
		os.Exit(1)
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		_, errr := udpConn.Write([]byte(line))
		if errr != nil {
			fmt.Fprintf(os.Stderr, "Error writing line: %s : %v\n", line, err)
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", line)
	}

}
