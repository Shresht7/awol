package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

// ----
// MAIN
// ----

// The main entrypoint of the application
func main() {

	// Parse the command-line arguments
	args := parseCommandLineArgs()

	// Check for help flag and show the help message
	if args.Help {
		helpMessage()
		return
	}

	// Load the configuration file
	cfgPath := getConfigPath()
	config, err := readConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		return
	}

	// The argument is expected to be a MAC address or an alias defined in the config
	macAddress := args.Mac
	macAlias := ""

	// Check if the provided MAC address is an alias in the config
	if val, exists := config.Aliases[strings.ToLower(args.Mac)]; exists {
		macAddress = val
		macAlias = args.Mac
	}

	// Parse the MAC address
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing MAC Addresss [%s]: %v\n", macAddress, err)
		return
	}

	// Create the magic packet
	magicPacket := makeMagicPacket(mac)

	// Send the magic packet via UDP broadcast (standard port for Wake-on-LAN is 9)
	conn, err := net.Dial("udp", fmt.Sprintf("255.255.255.255:%d", args.Port))
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

	if macAlias != "" {
		fmt.Printf("Magic packet sent to %s [%s]\n", macAlias, macAddress)
	} else {
		fmt.Printf("Magic packet sent to %s\n", macAddress)
	}
}

// ------------
// MAGIC PACKET
// ------------

// makeMagicPacket creates a Wake-on-LAN magic packet for the given hardware address.
// The magic packet consists of 6 bytes of 0xFF followed by 16 repetitions of the target MAC address.
// The resulting byte slice can be sent over the network to wake up the target device.
func makeMagicPacket(hardwareAddress net.HardwareAddr) []byte {
	var packet bytes.Buffer

	// 6 bytes of 0xFF
	packet.Write(bytes.Repeat([]byte{0xFF}, 6))

	// 16 repetitions of the MAC address
	for range 16 {
		packet.Write(hardwareAddress)
	}

	return packet.Bytes()
}

// -------------
// CONFIGURATION
// -------------

// Config struct to hold the aliases from the config file
type Config struct {
	Aliases map[string]string `json:"aliases"`
}

// getConfigPath retrieves the path to the configuration file from the environment variable "AWOL_CONFIG_PATH".
// If the environment variable is not set, it defaults to "~/.config/awol/config.json".
func getConfigPath() string {
	cfgPath := os.Getenv("AWOL_CONFIG_PATH")
	if cfgPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "~/.config/awol/config.json"
		}
		cfgPath = path.Join(home, ".config/awol/config.json")
	}
	return cfgPath
}

// Reads the configuration file from the specified path and returns the Config struct
func readConfig(cfgPath string) (Config, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening config file: %v\n", err)
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding config file: %v\n", err)
		return Config{}, err
	}

	return config, nil
}

// ------------
// COMMAND-LINE
// ------------

// The command-line arguments
type Args struct {
	Mac  string
	Help bool
	Port int
}

// Parse the command-line arguments and return an Args struct containing the parsed values
func parseCommandLineArgs() Args {
	help := flag.Bool("help", false, "Show help message")
	port := flag.Int("port", 9, "Port number to send the magic packet to")
	flag.Parse()

	return Args{
		Mac:  flag.Arg(0),
		Help: *help,
		Port: *port,
	}
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
