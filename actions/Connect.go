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
import "github.com/cookiengineer/hydra/receivers"
import "github.com/cookiengineer/hydra/types"

func Connect(host string, position string) error {

	screen, err0 := parsers.Xrandr()

	if err0 == nil {

		hostname, err1 := os.Hostname()

		if err1 == nil {

			ip := ""
			addrs, err2 := net.InterfaceAddrs()

			if err2 == nil {

				for _, addr := range addrs {

					if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

						if ipnet.IP.To4() != nil {
							ip = ipnet.IP.String()
							break
						}

					}

				}

			}

			if ip != "" {

				machine := types.Machine{
					Hostname: hostname,
					IP:       ip,
					Position: position,
					Screen:   screen,
				}

				data, err3 := json.Marshal(machine)

				if err3 == nil {

					url := fmt.Sprintf("http://%s:3000/connect", host)
					response, err4 := http.Post(url, "application/json", bytes.NewBuffer(data))

					if err4 == nil {

						defer response.Body.Close()

						if response.StatusCode == http.StatusOK {

							fmt.Println("Connected to hydra host:", host)

							return nil

						} else {
							return fmt.Errorf("connect request failed with status: %d", response.StatusCode)
						}

					} else {
						return err4
					}

				} else {
					return err3
				}

			} else {
				return errors.New("Could not determine local IP")
			}

		} else {
			return err1
		}

	} else {
		return err0
	}

}

func ReceiveEvents(host string) error {

	hostname, _ := os.Hostname()
	url := fmt.Sprintf("http://%s:3000/connect?hostname=%s", host, hostname)

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("receive events failed with status: %d", resp.StatusCode)
	}

	bridge, err0 := xorg.NewBridge(":0")

	if err0 != nil {
		return err0
	}

	defer bridge.Destroy()

	var virtual_screen *types.VirtualScreen = nil

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {

		line := scanner.Bytes()

		if len(line) == 0 || string(line) == "{}" {
			continue
		}

		var init_message struct {
			Type          string               `json:"type"`
			VirtualScreen *types.VirtualScreen `json:"virtual_screen"`
		}

		if err := json.Unmarshal(line, &init_message); err == nil && init_message.Type == "init" {

			virtual_screen = init_message.VirtualScreen

			fmt.Printf("Received virtual screen: %dx%d\n", virtual_screen.Width, virtual_screen.Height)

			continue

		}

		var mouse_event types.MouseEvent

		if err := json.Unmarshal(line, &mouse_event); err == nil && mouse_event.Type != 0 {

			receivers.ApplyMouseEvent(bridge, &mouse_event, virtual_screen, hostname)
			continue

		}

		var keyboard_event types.KeyboardEvent

		if err := json.Unmarshal(line, &keyboard_event); err == nil && keyboard_event.Type != 0 {

			receivers.ApplyKeyboardEvent(bridge, &keyboard_event)
			continue

		}

		fmt.Printf("Unknown event: %s\n", string(line))

	}

	return scanner.Err()

}
