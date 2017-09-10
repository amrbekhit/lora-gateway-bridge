package main

const configTemplate = `[general]
# ip:port to bind the UDP listener to
udp_bind = "{{ .General.UDPBind }}"

# debug=5, info=4, warning=3, error=2, fatal=1, panic=0
log_level = {{ .General.LogLevel }}

# skip the CRC status-check of received packets
skip_crc_check = {{ .General.SkipCRCCheck }}

[backend]
[backend.mqtt]
# mqtt server (e.g. scheme://host:port where scheme is tcp, ssl or ws)
server="{{ .Backend.MQTT.Server }}"

# mqtt server username (optional) 
username="{{ .Backend.MQTT.Username }}"

# mqtt server password (optional)
password="{{ .Backend.MQTT.Password }}"

# mqtt CA certificate file (optional)
ca_cert="{{ .Backend.MQTT.CACert }}"

# mqtt TLS certificate file (optional)
tls_cert="{{ .Backend.MQTT.TLSCert }}"

# mqtt TLS key file (optional)
tls_key="{{ .Backend.MQTT.TLSKey }}"
`
