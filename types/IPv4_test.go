package types

import "encoding/json"
import "testing"

func TestIPv4Parse(t *testing.T) {

	tests := []struct {
		input    string
		expected string
		valid    bool
	}{
		{"192.168.1.1", "192.168.1.1", true},
		{"10.0.0.1", "10.0.0.1", true},
		{"0.0.0.0", "0.0.0.0", true},
		{"255.255.255.255", "255.255.255.255", true},
		{"127.0.0.1", "127.0.0.1", true},
		{"invalid", "", false},
		{"192.168.1", "", false},
		{"256.256.256.256", "", false},
	}

	for _, test := range tests {

		ipv4 := ParseIPv4(test.input)

		if test.valid && ipv4 == nil {
			t.Errorf("Expected valid IPv4 for %s", test.input)
		}

		if !test.valid && ipv4 != nil {
			t.Errorf("Expected invalid IPv4 for %s", test.input)
		}

		if ipv4 != nil && ipv4.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, ipv4.String())
		}

	}

}

func TestIPv4IsIPv4(t *testing.T) {

	if !IsIPv4("192.168.1.1") {
		t.Error("Expected true for valid IPv4")
	}

	if IsIPv4("not-an-ip") {
		t.Error("Expected false for invalid string")
	}

	if IsIPv4("::1") {
		t.Error("Expected false for IPv6")
	}

}

func TestIPv4JSON(t *testing.T) {

	type Config struct {
		IP IPv4 `json:"ip"`
	}

	data := []byte(`{"ip":"192.168.1.1"}`)

	var cfg Config

	err := json.Unmarshal(data, &cfg)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if cfg.IP.String() != "192.168.1.1" {
		t.Errorf("Expected 192.168.1.1, got %s", cfg.IP.String())
	}

	marshalled, err := json.Marshal(cfg)

	if err != nil {
		t.Errorf("Expected no marshal error, got: %v", err)
	}

	if string(marshalled) != `{"ip":"192.168.1.1"}` {
		t.Errorf("Expected {\"ip\":\"192.168.1.1\"}, got %s", string(marshalled))
	}

}

