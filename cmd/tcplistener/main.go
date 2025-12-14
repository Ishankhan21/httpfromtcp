package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChannel := make(chan string)

	go func() {
		bytesLine := []byte{}
		for {
			bytes := make([]byte, 8)
			n, err := f.Read(bytes)
			if err == io.EOF {
				linesChannel <- string(bytesLine)
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				os.Exit(0)
			}
			data := bytes[:n]
			for _, v := range data {
				if string(v) != "\n" {
					bytesLine = append(bytesLine, v)
				}
				if string(v) == "\n" {
					linesChannel <- string(bytesLine)
					bytesLine = []byte{}
				}
			}
		}
		defer close(linesChannel)
	}()

	return linesChannel
}

func main() {

	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("read: %s\n", line)
		}
	}

}
