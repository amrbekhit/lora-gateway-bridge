/*
LoRa Gateway Bridge abstracts the packet_forwarder protocol into JSON over MQTT
	> documentation & support: https://docs.loraserver.io/lora-gateway-bridge
	> source & copyright information: https://github.com/brocaar/lora-gateway-bridge

Usage:
  lora-gateway-bridge [flags]
  lora-gateway-bridge [command]

<<<<<<< HEAD
COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --udp-bind value       ip:port to bind the UDP listener to (default: "0.0.0.0:1700") [$UDP_BIND]
   --mqtt-server value    mqtt server (e.g. scheme://host:port where scheme is tcp, ssl or ws) (default: "tcp://127.0.0.1:1883") [$MQTT_SERVER]
   --mqtt-username value  mqtt server username (optional) [$MQTT_USERNAME]
   --mqtt-password value  mqtt server password (optional) [$MQTT_PASSWORD]
   --mqtt-ca-cert value   mqtt CA certificate file (optional) [$MQTT_CA_CERT]
   --mqtt-tls-cert value  mqtt certificate file (optional) [$MQTT_TLS_CERT]
   --mqtt-tls-key value   mqtt key file of certificate (optional) [$MQTT_TLS_KEY]
   --skip-crc-check       skip the CRC status-check of received packets [$SKIP_CRC_CHECK]
   --log-level value      debug=5, info=4, warning=3, error=2, fatal=1, panic=0 (default: 4) [$LOG_LEVEL]
   --help, -h             show help
   --version, -v          print the version

COPYRIGHT:
   See http://github.com/brocaar/lora-gateway-bridge for copyright information

=======
Available Commands:
  help        Help about any command
  version     Print the LoRa Gateway Bridge version

Flags:
      --config string          config file (optional)
  -h, --help                   help for lora-gateway-bridge
      --log-level int          debug=5, info=4, error=2, fatal=1, panic=0 (default 4)
      --mqtt-ca-cert string    mqtt CA certificate file (optional)
      --mqtt-password string   mqtt server password (optional)
      --mqtt-server string     mqtt server (e.g. scheme://host:port where scheme is tcp, ssl or ws) (default "tcp://127.0.0.1:1883")
      --mqtt-username string   mqtt server username (optional)
      --skip-crc-check         skip the CRC status-check of received packets
      --udp-bind string        ip:port to bind the UDP listener to (default "0.0.0.0:1700")

Use "lora-gateway-bridge [command] --help" for more information about a command.

>>>>>>> Migrate to config file.
*/
package main
