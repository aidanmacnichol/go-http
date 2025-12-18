package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, 8)
	for {
		n, err := file.Read(buf)

		if n > 0 {
			fmt.Printf("read: %s\n", string(buf[:n]))
		}

		if err == io.EOF {
			return
		} else if err != nil {
			panic(err)
		}

	}

}
