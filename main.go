package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Check if a MAC address is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: No MAC address provided.")
		helpMessage()
		return
	}

	// Check for help flag and show the help message
	argument := os.Args[1]
	if argument == "-h" || argument == "--help" || argument == "help" {
		helpMessage()
		return
	}
	// otherwise, treat the argument as a MAC address
	macAddress := argument

	// Parse the MAC address
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing MAC Addresss [%s]: %v\n", macAddress, err)
		return
	}

	// Create the magic packet
	magicPacket := makeMagicPacket(mac)

	// Send the magic packet via UDP broadcast (standard port for Wake-on-LAN is 9)
	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error establishing UDP connection: %v\n", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(magicPacket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending magic packet: %v\n", err)
		return
	}

	fmt.Printf("Magic packet sent to %s\n", macAddress)
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

// Prints the help message for the command-line
func helpMessage() {
	help := strings.Builder{}
	help.WriteString("awol - a wake-on-lan utility\n\n")
	help.WriteString("Usage: awol <mac_address>\n\n")
	help.WriteString("Example:\n")
	help.WriteString("  awol A1:2B:C3:4D:5E:F7\n")
	fmt.Print(help.String())
}
