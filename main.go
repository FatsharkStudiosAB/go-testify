package main

import (
	"bytes"
	"flag"
	"go-testify/internal/stingray"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Not actually waiting for a signal, just polling the buffer for the signal string
func waitForSignal(buf *bytes.Buffer, signal string) {
	log.Printf("Waiting for signal: %s", signal)
	for {
		output := buf.String()
		if len(output) >= len(signal) && strings.Contains(output, signal) {
			log.Printf("Received signal: %s", signal)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func runTestCase(connector *stingray.Connector, testCase string, args ...string) {
	log.Printf("Running test case %s with args %v", testCase, args)
	messageType := "script"
	system := "Testify"
	stingray.ConsoleSend(connector, messageType, system, testCase, args...)
}

func main() {
	win32_data_dir := flag.String("data_dir", "E:/Projects/Bishop_data/win32", "Path to the win32 data directory")
	flag.Parse()

	log.Println("Initializing Testify in Go using Testify...")

	cmd := exec.Command(stingray.Exe_Directory + stingray.Exe_File)
	args := []string{"--data-dir", *win32_data_dir, "--disable-vsync", "--lua-discard-bytecode", "--port", stingray.Port, "--suppress-messagebox", "-game", "-testify", "-debug_testify", "-network_lan", "-skip_gamertag_popup", "-multiplayer_mode", "host"}
	// the below args are only temporary, they're the same args that are used for the spawn_all_enemies test
	args = append(args, strings.Split("-game -skip_first_character_creation -skip_prologue -skip_cinematics -mission spawn_all_enemies -dev -crash_on_account_login_error -character_profile_selector 1 -chunk_detector_free_flight_camera_raycast -chunk_lod_free_flight_camera_raycast -disable_pacing", " ")...)
	cmd.Args = append(cmd.Args, args...)

	var buf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalf("Error starting Stingray process: %v", err)
	}
	defer cmd.Process.Kill()

	connector := stingray.NewConnector()
	log.Printf("Server will start at %s:%s over %s\n", stingray.Address, stingray.Port, stingray.Protocol)
	timeout := 10 * time.Second
	maxRetries := 5
	err = connector.Connect(timeout, maxRetries)
	if err != nil {
		log.Fatalf("Error connecting: %v", err.Error())
	}
	defer connector.Disconnect()

	c := make(chan int)
	go func() {
		waitForSignal(&buf, "[Lua] INFO [Testify] Ready!")
		stingray.ConsoleSend(connector, "message", "Testify", "Hello World!")
		// Below is disabled until sending to console works properly
		// stingray.RunLuaFunction(connector, "Testify:ready_signal_received")
		// runTestCase(connector, "CombatTestCases.spawn_all_enemies", "{ kill_timer = 5 }")
	}()
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
