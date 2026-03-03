package actions

import "bufio"
import "bytes"
import "encoding/json"
import "errors"
import "fmt"
import "net"
import "net/http"
import "os"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/parsers"
import "github.com/cookiengineer/hydra/types"

// Connect sends the client machine info to the server
func Connect(host string, position string) error {

	screen, err := parsers.Xrandr()
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	ip := ""
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
					break
				}
			}
		}
	}

	if ip == "" {
		return errors.New("could not determine local IP")
	}

	machine := types.Machine{
		Hostname: hostname,
		IP:       ip,
		Position: position,
		Screen:   *screen,
	}

	data, err := json.Marshal(machine)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:3000/connect", host)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("connect request failed")
	}

	fmt.Println("Connected to hydra host:", host)
	return nil
}

func ReceiveEvents(host string, virtualScreen *types.VirtualScreen) error {

	hostname, _ := os.Hostname()
	url := fmt.Sprintf("http://%s:3000/connect?hostname=%s", host, hostname)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	bridge, err0 := xorg.NewBridge(":0")

	if err0 == nil {


	}

	// Attach virtual screen for offset calculations
	state.VirtualScreen = virtualScreen

	go bridge.Init()

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 || string(line) == "{}" {
			continue // keep-alive ping
		}

		// Try MouseEvent first
		var me types.MouseEvent
		if err := json.Unmarshal(line, &me); err == nil && me.Type != 0 {
			// Treat DX/DY as absolute global coordinates
			state.WarpPointer(me.DX, me.DY)
			// receivers.ApplyMouseEvent(state, &me)
			continue
		}

		// Try KeyboardEvent
		var ke types.KeyboardEvent
		if err := json.Unmarshal(line, &ke); err == nil && ke.Type != 0 {
			// receivers.ApplyKeyboardEvent(state, &ke)
			continue
		}
	}

	bridge.Destroy()

	return scanner.Err()

}

