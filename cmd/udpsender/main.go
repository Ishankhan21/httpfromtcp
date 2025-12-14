package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	// Line 13: Resolves the address "localhost:42069" into a UDP address structure
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	// Line 18: Creates a UDP "connection" (but it's not really a connection!)
	// This doesn't establish a persistent connection - it just sets up a socket
	// that knows where to send packets
	con, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	reader := bufio.NewReader(os.Stdin)
	// Line 24-37: Reads from stdin and sends each line as a UDP packet
	for {
		fmt.Println(">")

		stringLine, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		_, err = con.Write([]byte(stringLine))
		if err != nil {
			log.Fatal(err)
		}
	}

}
