package actions

import "bufio"
import "bytes"
import "encoding/json"
import "errors"
import "fmt"
import "io"
import "net"
import "net/http"
import "os"
import "strings"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/parsers"
import "github.com/cookiengineer/hydra/receivers"
import "github.com/cookiengineer/hydra/types"

func controllerName(host string) string {

	name := strings.TrimSpace(os.Getenv("HYDRA_CONTROLLER"))

	if name == "" {
		name = host
	}

	return strings.ToLower(name)

}

func localIP() (string, error) {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, addr := range addrs {

		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}

	}

	return "", errors.New("Could not determine local IP")

}

func sendHandoff(host string, hostname string, direction string) {

	handoff := types.HandoffEvent{
		Type:      "handoff",
		Machine:   hostname,
		Direction: direction,
	}

	data, err := json.Marshal(handoff)

	if err != nil {
		return
	}

	url := fmt.Sprintf("http://%s:%s/event", host, types.Port())

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Protocol", "hydra")
	request.Header.Set("X-Hydra-Controller", controllerName(host))

	client := &http.Client{}
	response, err := client.Do(request)

	if err == nil {
		response.Body.Close()
	}

}

func Connect(host string, position string) error {

	screen, err0 := parsers.Xrandr()

	if err0 != nil {
		return err0
	}

	hostname, err1 := os.Hostname()

	if err1 != nil {
		return err1
	}

	ip, err2 := localIP()

	if err2 != nil {
		return err2
	}

	machine := types.Machine{
		Hostname: hostname,
		IP:       ip,
		Position: position,
		Screen:   screen,
	}

	data, err3 := json.Marshal(machine)

	if err3 != nil {
		return err3
	}

	url := fmt.Sprintf("http://%s:%s/connect", host, types.Port())

	request, err4 := http.NewRequest("POST", url, bytes.NewBuffer(data))

	if err4 != nil {
		return err4
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Protocol", "hydra")
	request.Header.Set("X-Hydra-Controller", controllerName(host))

	client := &http.Client{}
	response, err5 := client.Do(request)

	if err5 != nil {
		return err5
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("connect request failed with status: %d", response.StatusCode)
	}

	fmt.Println("Connected to hydra host:", host)

	return receiveEvents(response.Body, host, hostname)

}

func receiveEvents(body io.Reader, host string, hostname string) error {

	bridge, err0 := xorg.NewBridge(defaultDisplay())

	if err0 != nil {
		return err0
	}

	defer bridge.Destroy()
	defer receivers.ReleaseAll(bridge)

	var virtual_screen *types.VirtualScreen = nil

	state := types.NewGlobalState()

	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {

		line := scanner.Bytes()

		if len(line) == 0 || string(line) == "{}" {
			continue
		}

		var envelope types.Event

		if err := json.Unmarshal(line, &envelope); err != nil {
			fmt.Printf("Unknown event: %s\n", string(line))
			continue
		}

		switch envelope.Type {

		case "init":

			var init_event types.InitEvent

			if err := envelope.Unmarshal(&init_event); err != nil {
				continue
			}

			virtual_screen = init_event.VirtualScreen

			for _, ws := range init_event.Workspaces {
				ws.Windows = []types.Window{}
				state.Workspaces[ws.Name] = &ws
			}

			state.SetActiveWorkspace(init_event.ActiveWorkspace)

			if virtual_screen != nil {
				fmt.Printf("Received virtual screen: %dx%d\n", virtual_screen.Width, virtual_screen.Height)
			}

		case "workspace":

			var workspace_event types.WorkspaceEvent

			if err := envelope.Unmarshal(&workspace_event); err == nil {
				receivers.ApplyWorkspaceEvent(bridge, state, &workspace_event)
			}

		case "focus":

			var focus_event types.FocusEvent

			if err := envelope.Unmarshal(&focus_event); err == nil {

				found := receivers.ApplyFocusEvent(bridge, &focus_event)

				if !found && (focus_event.Direction == "left" || focus_event.Direction == "right") {
					sendHandoff(host, hostname, focus_event.Direction)
				}

			}

		case "tile":

			var tile_event types.TileEvent

			if err := envelope.Unmarshal(&tile_event); err == nil {
				receivers.ApplyTileEvent(bridge, &tile_event, virtual_screen, hostname)
			}

		case "reset":

			var reset_event types.ResetEvent

			if err := envelope.Unmarshal(&reset_event); err == nil {
				receivers.ApplyResetEvent(bridge)
			}

		case "mouse":

			var mouse_event types.MouseEvent

			if err := envelope.Unmarshal(&mouse_event); err == nil {
				receivers.ApplyMouseEvent(bridge, &mouse_event, virtual_screen, hostname)
			}

		case "keyboard":

			var keyboard_event types.KeyboardEvent

			if err := envelope.Unmarshal(&keyboard_event); err == nil {
				receivers.ApplyKeyboardEvent(bridge, &keyboard_event)
			}

		default:

			fmt.Printf("Unknown event: %s\n", string(line))

		}

	}

	return scanner.Err()

}
