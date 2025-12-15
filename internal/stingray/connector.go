package stingray

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	Address       = "localhost"
	Exe_Directory = "C:/BitSquidBinaries/bishop/engine/win64/dev/"
	Exe_File      = "stingray_win64_dev_x64.exe"
	Port          = "14030"
	Protocol      = "tcp"
)

func NewProcess() *os.Process {
	command := Exe_File
	args := []string{"--disable-vsync", "--lua-discard-bytecode", "--port", "14030", "--suppress-messagebox", "-game", "-testify", "-debug_testify", "-network_lan", "-skip_gamertag_popup", "-multiplayer_mode", "host"}

	fmt.Printf("Launching Darktide executable %s in %s with args %s", Exe_File, Exe_Directory, args)

	procAttr := new(os.ProcAttr)
	procAttr.Dir = Exe_Directory
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	process, err := os.StartProcess(command, args, procAttr)
	if err != nil {
		fmt.Println("Error starting process:", err)
		panic(err)
	} else {
		fmt.Println("Process started with PID:", process.Pid)
		return process
	}
}

type Connector struct {
	address  string
	conn     net.Conn
	port     string
	protocol string
}

func NewConnector() *Connector {
	return &Connector{
		address:  Address,
		port:     Port,
		protocol: Protocol,
	}
}

func (connector *Connector) Connect(timeout time.Duration, maxRetries int) error {
	connector.conn = nil
	for retry := 0; retry < maxRetries; retry++ {
		conn, err := net.Dial(connector.protocol, connector.address+":"+connector.port)
		if err == nil {
			fmt.Println("Connected to game client...")
			defer conn.Close()
			connector.conn = conn
			return nil
		}
		fmt.Println("Failed to connect", err)
		time.Sleep(timeout)
	}
	return fmt.Errorf("failed to connect after %d retries", maxRetries)
}

func (connector *Connector) Disconnect() {
	connector.conn.Close()
}

// Placeholder for testing
func (connector *Connector) Hello() error {
	_, err := connector.conn.Write([]byte("Hello from Testify Client"))
	if err != nil {
		fmt.Println("Error writing to server:", err.Error())
	}

	return err
}
