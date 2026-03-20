package main

import (
	"bytes"
	"net"
)

// ------------
// MAGIC PACKET
// ------------

// makeMagicPacket creates a Wake-on-LAN magic packet for the given hardware address.
// The magic packet consists of 6 bytes of 0xFF followed by 16 repetitions of the target MAC address.
// The resulting byte slice can be sent over the network to wake up the target device.
func makeMagicPacket(hardwareAddress net.HardwareAddr) []byte {
	packet := make([]byte, 0, 6+16*len(hardwareAddress))

	// 6 bytes of 0xFF
	packet = append(packet, bytes.Repeat([]byte{0xFF}, 6)...)

	// 16 repetitions of the MAC address
	for range 16 {
		packet = append(packet, hardwareAddress...)
	}

	return packet
}
