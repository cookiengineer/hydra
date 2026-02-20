package actions

import "bytes"
import "encoding/json"
import "errors"
import "fmt"
import "net"
import "net/http"
import "os"
import "github.com/cookiengineer/hydra/parsers"
import "github.com/cookiengineer/hydra/types"

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

