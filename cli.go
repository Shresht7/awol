package main

import (
	"flag"
	"fmt"
	"strings"
)

// ----------------------
// COMMAND-LINE INTERFACE
// ----------------------

// The command-line arguments
type Args struct {
	SubCmd  string
	Mac     string
	Help    bool
	Version bool
	Port    int
	Rest    []string
}

// Parse the command-line arguments and return an Args struct containing the parsed values
func parseCommandLineArgs() Args {
	help := flag.Bool("help", false, "Show help message")
	version := flag.Bool("version", false, "Show version information")
	port := flag.Int("port", 0, "Port number to send the magic packet to")
	flag.Parse()

	mac := ""
	subcmd := flag.Arg(0)
	switch subcmd {
	case "help":
		*help = true
	case "wake", "call", "recall":
		mac = flag.Arg(1)
	default:
		mac = subcmd
	}

	return Args{
		SubCmd:  subcmd,
		Mac:     mac,
		Rest:    flag.Args()[1:],
		Help:    *help,
		Version: *version,
		Port:    *port,
	}
}

// Prints the help message for the command-line
func showHelp() {
	help := strings.Builder{}
	help.WriteString("awol - a wake-on-lan utility\n\n")
	help.WriteString("Usage: awol <mac>\n\n")
	help.WriteString("Commands:\n")
	help.WriteString("  wake <mac>\t\tSend a magic packet to the specified MAC address [aliases: call, recall]\n")
	help.WriteString("  list\t\t\tList all defined aliases in the config file\n")
	help.WriteString("  alias <alias> <mac>\tDefine a new alias for a MAC address in the config file\n")
	help.WriteString("  remove <alias>\tRemove an existing alias from the config file\n\n")
	help.WriteString("Flags:\n")
	help.WriteString("  --port <number>\tSpecify the port number to send the magic packet to (default: 9)\n")
	help.WriteString("  --version\t\tShow version information\n")
	help.WriteString("  --help\t\tShow this help message\n\n")
	help.WriteString("Example:\n")
	help.WriteString("  awol A1:2B:C3:4D:5E:F7\t# Send magic packet to the specified MAC address\n")
	help.WriteString("  awol wake skynet --port 7\t\t# Send magic packet to the specified MAC address using an alias on port 7\n")
	fmt.Print(help.String())
}

const version = "v0.2.0"

func showVersion() {
	fmt.Printf("%s\n", version)
}
