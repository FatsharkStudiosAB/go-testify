package stingray

import (
    "fmt"
    "net"
)

const (
    NetworkType = "tcp"
)

type Connector struct {
    address string
    conn    net.Conn
}

func NewConnector(address string) *Connector {
    return &Connector{
        address: address,
    }
}

func (connector *Connector) Connect() error {
    connector.conn = nil
    conn, err := net.Dial(NetworkType, connector.address)

    if err != nil {
        fmt.Println("Failed to connect", err)
        return err
    }

    defer conn.Close()
    connector.conn = conn
    return nil
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