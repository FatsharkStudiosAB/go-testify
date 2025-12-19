package main

import (
	"go-testify/internal/darktide"
	"go-testify/internal/stingray"
	"log"
	"time"
)

// type TestCase struct {
// 	name        string
// 	id          string
// 	script_func string
// 	script_args map[string]any
// }

// type Session struct {
// 	project_config_type string
// 	extra_session_args  string
// }

// type TestSuite struct {
// 	name       string
// 	id         string
// 	sessions   []Session
// 	test_cases []TestCase
// }

func main() {
	log.Println("Initializing Testify in Go using Testify...")

	client := darktide.NewDarktideClient()
	_, err := client.Start()
	if err != nil {
		log.Fatalf("Error starting Darktide client process: %v", err)
	}
	defer client.Stop()

	connector := stingray.NewConnector()
	log.Printf("Connecting to %s:%s over %s\n", stingray.Address, stingray.Port, stingray.Protocol)
	timeout := 1 * time.Second
	maxRetries := 20
	err = connector.Connect(timeout, maxRetries)
	if err != nil {
		log.Fatalf("Error connecting: %v", err.Error())
	}
	defer connector.Disconnect()

	c := make(chan int)
	go func() {
		time.Sleep(10 * time.Second) // Wait a bit for the game to initialize
		//client.WaitForLuaReadySignal(&buf)
		stingray.ConsoleSend(connector, "script", "Testify", "Testify:ready_signal_received", map[string]any{})
		stingray.ConsoleSend(connector, "script", "Testify", "CombatTestCases.spawn_all_enemies", map[string]any{"kill_timer": 10})
	}()
	go func() {
		err := client.Wait()
		if err != nil {
			log.Printf("Stingray process exited with error: %v", err)
		} else {
			log.Println("Stingray process exited successfully.")
		}
		c <- 1
	}()
	<-c

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
