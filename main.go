package main

import (
	"fmt"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

func main() {
	fmt.Println("Initializing Testify in Go using Testify...")

	fmt.Printf("Server will start at %s:%s over %s\n", SERVER_HOST, SERVER_PORT, SERVER_TYPE)
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		panic(err)
	}

	_, err = connection.Write([]byte("Hello from Testify Client"))
	if err != nil {
		fmt.Println("Error writing to server:", err.Error())
	}

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err.Error())
	}
	fmt.Println("Message from Server: " + string(buffer[:mLen]))
	defer connection.Close()
}
