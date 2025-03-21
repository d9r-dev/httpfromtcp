package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		var currentLine string
		for {
			buffer := make([]byte, 8)
			_, err := f.Read(buffer)
			if err != nil {
				break
			}

			parts := bytes.Split(buffer, []byte("\n"))
			if len(parts) == 1 {
				currentLine = currentLine + string(parts[0][:])
			} else {
				for i := 0; i < len(parts); i++ {
					if i != len(parts)-1 {
						currentLine = currentLine + string(parts[i][:])
						lines <- currentLine
						currentLine = ""
					} else {
						currentLine = currentLine + string(parts[i][:])
					}
				}

			}
		}
		if len(currentLine) != 0 {
			lines <- currentLine
		}
		err := f.Close()
		if err != nil {
			log.Fatalf("Could not close file.")
		}
		close(lines)
	}()

	return lines
}
