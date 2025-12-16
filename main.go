package main

import (
	"bytes"
	"flag"
	"go-testify/internal/stingray"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

func waitForSignal(stdout io.ReadCloser, signal string) {
	buf := new(bytes.Buffer)
	go func() {
		_, err := io.Copy(buf, stdout)
		if err != nil {
			log.Printf("Error copying stdout: %v", err)
		}
	}()
	log.Printf("Waiting for signal: %s", signal)
	for {
		output := buf.String()
		log.Printf("Output: %s", output)
		if len(output) >= len(signal) && strings.Contains(output, signal) {
			log.Printf("Received signal: %s", signal)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	win32_data_dir := flag.String("data_dir", "E:/Projects/Bishop_data/win32", "Path to the win32 data directory")
	flag.Parse()

	log.Println("Initializing Testify in Go using Testify...")

	cmd := exec.Command(stingray.Exe_Directory + stingray.Exe_File)
	args := []string{"--data-dir", *win32_data_dir, "--disable-vsync", "--lua-discard-bytecode", "--port", stingray.Port, "--suppress-messagebox", "-game", "-testify", "-debug_testify", "-network_lan", "-skip_gamertag_popup", "-multiplayer_mode", "host"}
	cmd.Args = append(cmd.Args, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating StdoutPipe for Stingray process: %v", err)
	}

	err = cmd.Start()
	defer cmd.Process.Kill()
	if err != nil {
		log.Fatalf("Error starting Stingray process: %v", err)
	}

	c := make(chan int)
	go waitForSignal(stdout, "[Lua] INFO [Testify] Ready!")
	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Printf("Stingray process exited with error: %v", err)
		} else {
			log.Println("Stingray process exited successfully.")
		}
		c <- 1
	}()
	<-c

	// fmt.Printf("Server will start at %s:%s over %s\n", stingray.Address, stingray.Port, stingray.Protocol)
	// connector := stingray.NewConnector()
	// timeout := 10 * time.Second
	// maxRetries := 5
	// err = connector.Connect(timeout, maxRetries)
	// if err != nil {
	// 	fmt.Println("Error connecting:", err.Error())
	// 	panic(err)
	// }

	// err = connector.Hello()
	// if err != nil {
	// 	fmt.Println("Error writing to server:", err.Error())
	// }

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
