---
title: Configuration
menu:
    main:
        parent: install
        weight: 5
---

## Configuration

### Gateway

Modify the [packet-forwarder](https://github.com/lora-net/packet_forwarder)
of your gateway so that it will send its data to the LoRa Gateway Bridge.
You will need to change the following configuration keys:

* `server_address` to the IP address / hostname of the LoRa Gateway Bridge
* `serv_port_up` to `1700` (the default port that LoRa Gateway Bridge is using)
* `serv_port_down` to `1700` (same)

### LoRa Gateway Bridge

The `lora-gateway-bridge` has the following command-line flags:

```text
LoRa Gateway Bridge abstracts the packet_forwarder protocol into JSON over MQTT
        > documentation & support: https://docs.loraserver.io/lora-gateway-bridge
        > source & copyright information: https://github.com/brocaar/lora-gateway-bridge

Usage:
  lora-gateway-bridge [flags]
  lora-gateway-bridge [command]

Available Commands:
  configfile  Print the LoRa Gateway configuration file
  help        Help about any command
  version     Print the LoRa Gateway Bridge version

Flags:
  -c, --config string   path to configuration file (optional)
  -h, --help            help for lora-gateway-bridge
      --log-level int   debug=5, info=4, error=2, fatal=1, panic=0 (default 4)

Use "lora-gateway-bridge [command] --help" for more information about a command.
```

#### Configuration file

By default `lora-gateway-bridge` will look in the following order for a
configuration at the following paths when `--config` / `-c` is unset:

* `lora-gateway-bridge.toml`
* `$HOME/.config/lora-gateway-bridge/lora-gateway-bridge.toml`
* `/etc/lora-gateway-bridge/lora-gateway-bridge.toml`

To load configuration from a different location, set the `--config` / `-c`
flag, or use the `CONFIG` environment variable.

Example configuration file:

```toml
[general]
# ip:port to bind the UDP listener to
#
# Example: 0.0.0.0:1700 to listen on port 1700 on all network interfaces.
# This is the listeren to which the packet-forwarder forwards its data.
udp_bind = "0.0.0.0:1700"

# debug=5, info=4, warning=3, error=2, fatal=1, panic=0
log_level = 4

# Skip the CRC status-check of received packets
#
# This is only has effect when the packet-forwarder is configured to forward
# LoRa frames with CRC errors.
skip_crc_check = false

# Configuration for the MQTT backend.
[backend.mqtt]
# MQTT server (e.g. scheme://host:port where scheme is tcp, ssl or ws)
server="tcp://127.0.0.1:1883"

# Connect with the given username (optional)
username=""

# Connect with the given password (optional)
password=""

# CA certificate file (optional)
#
# Use this when setting up a secure connection (when server uses ssl://...)
# but the certificate used by the server is not trusted by any CA certificate
# on the server (e.g. when self generated).
ca_cert=""

# mqtt TLS certificate file (optional)
tls_cert=""

# mqtt TLS key file (optional)
tls_key=""
```

#### Warning: deprecation warning! update your configuration

When you see this warning, you need to update your configuration!
Before LoRa Gateway Bridge 2.3.0 environment variables were used for setting
configuration flags. Since LoRa Gateway Bridge 2.3.0 the configuration format
has changed.

When installed from a `.deb` package, this is the recommended way to upgrade:

```bash
# load the environment variables
set -a
source /etc/default/lora-gateway-bridge

# create the configuration directory
mkdir /etc/lora-gateway-bridge

# generate new configuration file, pre-filled with the configuration set
# through the environment variables
lora-gateway-bridge configfile > /etc/lora-gateway-bridge/lora-gateway-bridge.toml

# "remove" the old configuration
mv /etc/default/lora-gateway-bridge /etc/default/lora-gateway-bridge.old
```
