package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	// Generate some random data
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, 16)
	rand.Read(data)

	// Send the random data to the server
	_, err = conn.Write(data)
	if err != nil {
		log.Fatalf("Failed to send data: %v", err)
	}

	fmt.Printf("Sent: %x\n", data)
}
