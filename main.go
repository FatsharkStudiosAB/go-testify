package main

import (
	"fmt"
	"go-testify/internal/stingray"
	"time"
)

func main() {
	fmt.Println("Initializing Testify in Go using Testify...")

	process := stingray.NewProcess()
	defer process.Kill()

	fmt.Printf("Server will start at %s:%s over %s\n", stingray.Address, stingray.Port, stingray.Protocol)
	connector := stingray.NewConnector()
	timeout := 10 * time.Second
	maxRetries := 5
	err := connector.Connect(timeout, maxRetries)
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
