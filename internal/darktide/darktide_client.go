package darktide

import (
	"bytes"
	"flag"
	"go-testify/internal/shell"
	"go-testify/internal/stingray"
	"log"
	"os/exec"
	"strings"
	"time"
)

type DarktideClient struct {
	command *exec.Cmd
}

func NewDarktideClient() *DarktideClient {
	// D:/svn/bishop/trunk_data/win32
	win32_data_dir := flag.String("data_dir", "E:/Projects/Bishop_data/win32", "Path to the win32 data directory")
	flag.Parse()
	cmd := exec.Command(stingray.Exe_Directory + stingray.Exe_File)
	args := []string{"--data-dir", *win32_data_dir, "--disable-vsync", "--lua-discard-bytecode", "--port", stingray.Port, "--suppress-messagebox", "-game", "-testify", "-debug_testify", "-network_lan", "-skip_gamertag_popup", "-multiplayer_mode", "host"}
	// the below args are only temporary, they're the same args that are used for the spawn_all_enemies test
	args = append(args, strings.Split("-game -skip_first_character_creation -skip_prologue -skip_cinematics -mission spawn_all_enemies -dev -crash_on_account_login_error -character_profile_selector 1 -chunk_detector_free_flight_camera_raycast -chunk_lod_free_flight_camera_raycast -disable_pacing", " ")...)
	cmd.Args = append(cmd.Args, args...)

	return &DarktideClient{command: cmd}
}

// TODO parameterize this
func (client DarktideClient) Start() (bytes.Buffer, error) {
	localShell := shell.NewLocalShell()
	buf, _, err := localShell.ExecuteCommand(client.command)
	if err != nil {
		log.Fatalf("Error starting Stingray process: %v", err)
	}

	return buf, err
}

func (client DarktideClient) Stop() {
	client.command.Process.Kill()
}

func (client DarktideClient) Wait() error {
	return shell.WaitForCommand(client.command)
}

func (client DarktideClient) WaitForLuaReadySignal(buf *bytes.Buffer) {
	waitForSignal(buf, "[Lua] INFO [Testify] Ready!")
}

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
