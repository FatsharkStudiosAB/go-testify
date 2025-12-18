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
}

func NewDarktideClient() *DarktideClient {
	return &DarktideClient{}
}

// TODO parameterize this
func (client DarktideClient) Start() (bytes.Buffer, error) {
	win32_data_dir := flag.String("data_dir", "E:/Projects/Bishop_data/win32", "Path to the win32 data directory")
	flag.Parse()
	cmd := exec.Command(stingray.Exe_Directory + stingray.Exe_File)
	args := []string{"--data-dir", *win32_data_dir, "--disable-vsync", "--lua-discard-bytecode", "--port", stingray.Port, "--suppress-messagebox", "-game", "-testify", "-debug_testify", "-network_lan", "-skip_gamertag_popup", "-multiplayer_mode", "host"}
	// the below args are only temporary, they're the same args that are used for the spawn_all_enemies test
	args = append(args, strings.Split("-game -skip_first_character_creation -skip_prologue -skip_cinematics -mission spawn_all_enemies -dev -crash_on_account_login_error -character_profile_selector 1 -chunk_detector_free_flight_camera_raycast -chunk_lod_free_flight_camera_raycast -disable_pacing", " ")...)
	cmd.Args = append(cmd.Args, args...)

	localShell := shell.NewLocalShell()
	buf, _, err := localShell.ExecuteCommand(cmd)
	if err != nil {
		log.Fatalf("Error starting Stingray process: %v", err)
	}
	defer cmd.Process.Kill()
/* TODO makee this work :)
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
*/

	return buf, err
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