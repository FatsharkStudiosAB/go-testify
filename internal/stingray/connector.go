package stingray

import (
	"fmt"
	"net"
	"time"
)

const (
	Address       = "localhost"
	Exe_Directory = "C:/BitSquidBinaries/bishop/engine/win64/dev/"
	Exe_File      = "stingray_win64_dev_x64.exe"
	Port          = "14030"
	Protocol      = "tcp"
)

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
