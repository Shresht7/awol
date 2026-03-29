package main

import (
	"bytes"
	"fmt"
	"net"
	"time"
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
func broadcastMagicPacket(source, network string, magicPacket []byte) error {
	// Validate the source IP address if provided
	sourceIP := net.ParseIP(source)
	if source != "" && sourceIP == nil {
		return fmt.Errorf("Invalid source IP address: %s", source)
	}

	// If a source IP address is provided, use it; otherwise, let the system choose the source IP
	sourceAddress := &net.UDPAddr{IP: sourceIP, Port: 0}

	// Resolve the broadcast address
	broadcastAddress, err := net.ResolveUDPAddr("udp4", network)
	if err != nil {
		return fmt.Errorf("Invalid broadcast address: %s", network)
	}

	// Establish a UDP connection to the broadcast address
	conn, err := net.DialUDP("udp", sourceAddress, broadcastAddress)
	if err != nil {
		return fmt.Errorf("Error establishing UDP connection: %w", err)
	}
	defer conn.Close() // Ensure the connection is closed when we're done, even if an error occurs

	// Set a write deadline to prevent hanging indefinitely if the network is unreachable
	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))

	// Write the magic packet to the UDP connection
	n, err := conn.Write(magicPacket)
	if err != nil {
		return fmt.Errorf("Error sending magic packet: %w", err)
	}
	if n != len(magicPacket) {
		return fmt.Errorf("Partial write: only %d of %d bytes sent", n, len(magicPacket))
	}

	return nil
}
