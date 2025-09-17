package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				if errors.Is(err, io.EOF) {
					log.Println("Server closed connection")
					os.Exit(0)
					return
				}

				log.Printf("Error reading from server: %v", err)
				return
			}
			if n > 0 {
				log.Printf("Server: %s", string(buf[:n]))
			}
		}
	}()

	fmt.Println("Connected to server, type your message and press enter, you will receive the messages from the server")

	for {
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}
		_, err = conn.Write([]byte(input + "\n"))
		if err != nil {
			log.Printf("Error sending to server: %v", err)
			break
		}
	}
}
