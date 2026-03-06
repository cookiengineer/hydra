package types

import "errors"
import "strings"

type Machine struct {
	Hostname string      `json:"hostname"`
	IP       string      `json:"ip"`
	Position string      `json:"position"` // left-of, right-of, above, below
	Screen   *Screen     `json:"screen"`
	Socket   chan []byte `json:"-"`
}

func (machine *Machine) Parse() error {

	machine.Hostname = strings.ToLower(strings.TrimSpace(machine.Hostname))
	machine.IP       = strings.ToLower(strings.TrimSpace(machine.IP))
	machine.Position = strings.ToLower(strings.TrimSpace(machine.Position))

	if machine.Hostname == "" {
		return errors.New("Invalid Hostname")
	}

	if machine.IP != "" {

		if IsIPv4(machine.IP) == true {

			ipv4 := ParseIPv4(machine.IP)

			if ipv4 != nil {
				machine.IP = ipv4.String()
			} else {
				return errors.New("Invalid IPv4")
			}

		} else if IsIPv6(machine.IP) == true {

			ipv6 := ParseIPv6(machine.IP)

			if ipv6 != nil {
				machine.IP = ipv6.String()
			} else {
				return errors.New("Invalid IPv6")
			}

		}

	} else {
		return errors.New("Invalid IP")
	}

	if machine.Position == "left-of" ||
		machine.Position == "right-of" ||
		machine.Position == "above" ||
		machine.Position == "below" ||
		machine.Position == "center" {
		// Do Nothing
	} else {
		return errors.New("Invalid Position")
	}

	if machine.Screen != nil {

		if machine.Screen.Width > 0 && machine.Screen.Height > 0 && len(machine.Screen.Monitors) > 0 {
			// Do Nothing
		} else {
			return errors.New("Invalid Screen Dimensions")
		}

	} else {
		return errors.New("Invalid Screen")
	}

	machine.Socket = make(chan []byte, 128)

	return nil

}

