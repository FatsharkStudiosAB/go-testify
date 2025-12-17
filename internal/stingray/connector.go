package stingray

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
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

type JsonStruct struct {
	Type   string `json:"type"`
	Script string `json:"script"`
}

func RunLuaFunction(connector *Connector, functionName string, args ...string) {
	function := functionName + "(" + strings.Join(args, ",") + ")"
	jsonData := JsonStruct{Type: "script", Script: "Application.console_send({ type = 'script', level = 'info', system = 'Testify', message = '" + function + "' })"}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	log.Printf("Marshaled JSON: %s", string(jsonBytes))
	// Create the message with length prefix
	// First 4 bytes: 0, Next 4 bytes: length of JSON
	msg := make([]byte, 8+len(jsonBytes))
	binary.BigEndian.PutUint32(msg[0:4], 0)
	binary.BigEndian.PutUint32(msg[4:8], uint32(len(jsonBytes)))
	copy(msg[8:], jsonBytes)
	log.Printf("Binary packed json: %s", msg)
	err = connector.Send(msg)
	if err != nil {
		log.Fatalf("Error sending JSON: %v", err)
	}
}

func (connector *Connector) Send(js []byte) error {
	_, err := connector.conn.Write(js)
	return err
}
