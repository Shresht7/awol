package main

import (
	"bytes"
	"fmt"
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

// broadcastMagicPacket sends the given payload to the specified network address using UDP.
func broadcastMagicPacket(network string, magicPacket []byte) error {
	conn, err := net.Dial("udp", network)
	if err != nil {
		return fmt.Errorf("Error establishing UDP connection: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(magicPacket)
	if err != nil {
		return fmt.Errorf("Error sending magic packet: %w", err)
	}
	return nil
}
