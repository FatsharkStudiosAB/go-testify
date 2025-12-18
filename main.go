package main

import (
	"go-testify/internal/darktide"
	"go-testify/internal/stingray"
	"log"
	"time"
)

func runTestCase(connector *stingray.Connector, testCase string, args ...string) {
	log.Printf("Running test case %s with args %v", testCase, args)
	messageType := "script"
	system := "Testify"
	stingray.ConsoleSend(connector, messageType, system, testCase, args...)
}

func main() {
	log.Println("Initializing Testify in Go using Testify...")

	client := darktide.NewDarktideClient()
	_, err := client.Start()
	if err != nil {
		log.Fatalf("Error starting Darktide client process: %v", err)
	}

	connector := stingray.NewConnector()
	log.Printf("Server will start at %s:%s over %s\n", stingray.Address, stingray.Port, stingray.Protocol)
	timeout := 10 * time.Second
	maxRetries := 5
	err = connector.Connect(timeout, maxRetries)
	if err != nil {
		log.Fatalf("Error connecting: %v", err.Error())
	}
	defer connector.Disconnect()

	/*  Might come in handy for actually reading signals from the server later
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err.Error())
	}
	fmt.Println("Message from Server: " + string(buffer[:mLen]))
	defer connection.Close()
	*/
}
