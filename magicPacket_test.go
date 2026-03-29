package main

import (
	"bytes"
	"net"
	"testing"
)

func TestMakeMagicPacket(t *testing.T) {
	testMACs := []string{
		"01:23:45:67:89:AB",
		"FF:FF:FF:FF:FF:FF",
		"00:00:00:00:00:00",
		"12-34-56-78-9A-BC",
	}

	for _, macStr := range testMACs {
		// Parse the MAC address string into a net.HardwareAddr type.
		mac, err := net.ParseMAC(macStr)
		if err != nil {
			t.Fatalf("Failed to parse MAC address %s: %v", macStr, err)
		}

		// Construct the expected magic packet for this MAC address.
		expected := bytes.Repeat([]byte{0xFF}, 6)
		for range 16 {
			expected = append(expected, mac...)
		}

		// Check that the generated magic packet matches the expected value.
		packet := makeMagicPacket(mac)
		if !bytes.Equal(packet, expected) {
			t.Errorf("Magic packet does not match expected value for MAC %s.\nGot: % X\nExpected: % X", macStr, packet, expected)
		}
	}
}
