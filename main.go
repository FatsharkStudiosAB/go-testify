package main

import (
	"fmt"
	"go-testify/internal/stingray"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
)

func main() {
	fmt.Println("Initializing Testify in Go using Testify...")

	fmt.Printf("Server will start at %s:%s over %s\n", SERVER_HOST, SERVER_PORT, stingray.NetworkType)
	connector := stingray.NewConnector(SERVER_HOST)
	err := connector.Connect()
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		panic(err)
	}

	err = connector.Hello()
	if err != nil {
		fmt.Println("Error writing to server:", err.Error())
	}

	/*
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err.Error())
	}
	fmt.Println("Message from Server: " + string(buffer[:mLen]))
	defer connection.Close()
	*/
}
