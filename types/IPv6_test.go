
package types

import "testing"

func TestIPv6Parse(t *testing.T) {

	tests := []struct {
		input    string
		expected string
		valid    bool
	}{
		{"::1", "0000:0000:0000:0000:0000:0000:0000:0001", true},
		{"[::1]", "0000:0000:0000:0000:0000:0000:0000:0001", true},
		{"fe80::1", "fe80:0000:0000:0000:0000:0000:0000:0001", true},
		{"invalid", "", false},
	}

	for _, test := range tests {

		ipv6 := ParseIPv6(test.input)

		if test.valid && ipv6 == nil {
			t.Errorf("Expected valid IPv6 for %s", test.input)
		}

		if !test.valid && ipv6 != nil {
			t.Errorf("Expected invalid IPv6 for %s", test.input)
		}

		if ipv6 != nil && ipv6.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, ipv6.String())
		}

	}

}
