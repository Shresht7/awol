package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wake <mac_address>")
		return
	}

	macAddress := os.Args[1]

	// Parse the MAC address
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing MAC Addresss [%s]: %v\n", macAddress, err)
		return
	}

	// Create the magic packet
	magicPacket := makeMagicPacket(mac)

	fmt.Printf("Magic Packet: %x\n", magicPacket)
}

// makeMagicPacket creates a Wake-on-LAN magic packet for the given hardware address.
// The magic packet consists of 6 bytes of 0xFF followed by 16 repetitions of the target MAC address.
// The resulting byte slice can be sent over the network to wake up the target device.
func makeMagicPacket(hardwareAddress net.HardwareAddr) []byte {
	var packet bytes.Buffer

	// 6 bytes of 0xFF
	packet.Write(bytes.Repeat([]byte{0xFF}, 6))

	// 16 repetitions of the MAC address
	for i := 0; i < 16; i++ {
		packet.Write(hardwareAddress)
	}

	return packet.Bytes()
}
