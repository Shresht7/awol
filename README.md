# `awol`

A simple command-line-interface to dispatch Wake-on-LAN magic packets!

This allows you to ***remotely*** switch on a computer/machine on your local network, (as long as it supports Wake-on-LAN and is configured to do so).

## Usage

Use `--help` or `-h` to see the help message.

To wake a computer, provide its MAC address as an argument to the `awol` command. For example:

```bash
awol <mac>
```
Example:

```bash
awol A1:2B:C3:4D:5E:F7 # Send magic packet to the specified MAC address
awol wake skynet # Send magic packet to the specified MAC address using an alias
```

### Options

- `--port <number>`: Specify the port number to send the magic packet to (default: 9)

### Configuration

You can define aliases for MAC addresses in a config file, which is stored at `~/.config/awol/config.json`. This allows you to use easy-to-remember names instead of typing out the full MAC address each time.

```json
{
  "aliases": {
    "skynet": "A1:2B:C3:4D:5E:F7",
    "hal": "B2:3C:D4:5E:F6:A8"
  }
}
```

You can use the `list` command to see all defined aliases, the `alias` command to add a new alias, and the `remove` command to delete an existing alias. 

For example:

```bash
awol list # List all defined aliases
awol alias skynet A1:2B:C3:4D:5E:F7 # Define a new alias 'skynet' for the specified MAC address
awol remove skynet # Remove the alias 'skynet'
```

### Configuration file

The configuration file is stored at `~/.config/awol/config.json` by default. You can override this path by setting the environment variable `AWOL_CONFIG_PATH` to a custom file path.

Config file fields:

- `broadcast` or `BroadcastAddress`: The broadcast IP address used to send magic packets. Defaults to `255.255.255.255` if not set.
- `port`: The UDP port to send the magic packet to. Defaults to `9` if not set.
- `aliases`: A map of friendly names to MAC addresses. Alias keys are treated case-insensitively and are normalized to lowercase; stored MAC addresses will be normalized to uppercase.

Example `config.json` with broadcast address and port:

```json
{
  "broadcast": "192.168.1.255",
  "port": 9,
  "aliases": {
    "skynet": "A1:2B:C3:4D:5E:F7",
    "hal": "B2:3C:D4:5E:F6:A8"
  }
}
```

Notes:

- If the configuration file does not exist, the application will use sensible defaults (`broadcast`: `255.255.255.255`, `port`: `9`). All fields are optional.
- Use the `alias` and `remove` commands to manage entries in the config file; the CLI will create the config directory/file as needed.

## Installation

```bash
go install github.com/Shresht7/awol@latest
```

## Wake-on-LAN

Wake-on-LAN (WoL) is a networking standard that allows a computer to be turned on or awakened from a low power state remotely. This is achieved by sending a specially crafted network packet, known as a _"magic packet"_ to the target computer's network interface. The magic packet contains the target computer's MAC address repeated multiple (16) times, which allows the network interface to recognize it and trigger the power-on process.

See [Wake-on-LAN - Wikipedia](https://en.wikipedia.org/wiki/Wake-on-LAN) for more details.

### Configuration

> [!IMPORTANT]
> To use Wake-on-LAN, the target computer must support it and be configured to allow it. This typically involves enabling WoL in the computer's **BIOS/UEFI** settings and ensuring that the network interface is set to allow wake-up events. Additionally, the computer must be connected to a power source and the network interface must be active (not in a completely powered-off state) for WoL to work.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
