package stingray

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	Address       = "localhost"
	Exe_Directory = "C:/BitSquidBinaries/bishop/engine/win64/dev/"
	Exe_File      = "stingray_win64_dev_x64.exe"
	PLAIN_JSON    = 0
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
			connector.conn = conn
			return nil
		}
		fmt.Println("Failed to connect", err)
		time.Sleep(timeout)
	}
	return fmt.Errorf("failed to connect after %d retries", maxRetries)
}

func (connector *Connector) Disconnect() {
	log.Println("Disconnecting from game client...")
	err := connector.conn.Close()
	if err != nil {
		log.Printf("Error disconnecting: %v", err)
	}
}

func send(conn net.Conn, data Message) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	log.Printf("Sending json: %s", jsonBytes)

	// Create message: 4 bytes type + 4 bytes length + JSON
	msg := make([]byte, 8+len(jsonBytes))
	binary.BigEndian.PutUint32(msg[0:4], PLAIN_JSON)
	binary.BigEndian.PutUint32(msg[4:8], uint32(len(jsonBytes)))
	copy(msg[8:], jsonBytes)

	_, err = conn.Write(msg)
	return err
}

type Message struct {
	Type   string `json:"type"`
	Script string `json:"script"`
}

func ConsoleSend(connector *Connector, messageType string, system string, message string, args map[string]any) {
	if len(args) > 0 {
		jsonArgs, err := json.Marshal(args)
		if err != nil {
			log.Fatalf("Error marshalling args: %v", err)
		}
		message = message + "('" + string(jsonArgs) + "')"
	}
	log.Printf("ConsoleSend called with type: %s, system: %s, message: %s", messageType, system, message)
	data := Message{
		Type:   "script",
		Script: "Application.console_send({ type = '" + messageType + "', level = 'info', system = '" + system + "', message = " + message + " })",
	}
	if err := send(connector.conn, data); err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
}
