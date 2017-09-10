package main

//go:generate ./doc.sh

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/brocaar/lora-gateway-bridge/backend/mqttpubsub"
	"github.com/brocaar/lora-gateway-bridge/gateway"
	"github.com/brocaar/lorawan"
)

type config struct {
	General struct {
		UDPBind      string `mapstructure:"udp_bind"`
		LogLevel     int    `mapstructure:"log_level"`
		SkipCRCCheck bool   `mapstructure:"skip_crc_check"`
	}

	Backend struct {
		MQTT struct {
			Server   string
			Username string
			Password string
			CACert   string `mapstructure:"ca_cert"`
			TLSCert  string `mapstructure:"tls_cert"`
			TLSKey   string `mapstructure:"tls_key"`
		}
	}
}

// C holds the global configuration
var C config

var version string // set by the compiler
var cfgFile string // config file

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (optional)")
	rootCmd.PersistentFlags().String("udp-bind", "0.0.0.0:1700", "ip:port to bind the UDP listener to")
	rootCmd.PersistentFlags().String("mqtt-server", "tcp://127.0.0.1:1883", "mqtt server (e.g. scheme://host:port where scheme is tcp, ssl or ws)")
	rootCmd.PersistentFlags().String("mqtt-username", "", "mqtt server username (optional)")
	rootCmd.PersistentFlags().String("mqtt-password", "", "mqtt server password (optional)")
	rootCmd.PersistentFlags().String("mqtt-ca-cert", "", "mqtt CA certificate file (optional)")
	rootCmd.PersistentFlags().String("mqtt-tls-cert", "", "")
	rootCmd.PersistentFlags().String("mqtt-tls-key", "", "")
	rootCmd.PersistentFlags().Bool("skip-crc-check", false, "skip the CRC status-check of received packets")
	rootCmd.PersistentFlags().Int("log-level", 4, "debug=5, info=4, error=2, fatal=1, panic=0")

	// for backwards compatibility
	viper.BindEnv("general.udp_bind", "UDP_BIND")
	viper.BindEnv("general.log_level", "LOG_LEVEL")
	viper.BindEnv("general.skip_crc_check", "SKIP_CRC_CHECK")
	viper.BindEnv("backend.mqtt.server", "MQTT_SERVER")
	viper.BindEnv("backend.mqtt.username", "MQTT_USERNAME")
	viper.BindEnv("backend.mqtt.password", "MQTT_PASSWORD")
	viper.BindEnv("backend.mqtt.ca_cert", "MQTT_CA_CERT")
	viper.BindEnv("backend.mqtt.tls_cert", "MQTT_TLS_CERT")
	viper.BindEnv("backend.mqtt.tls_key", "MQTT_TLS_KEY")

	// for backwards compatibility
	viper.BindPFlag("general.udp_bind", rootCmd.PersistentFlags().Lookup("udp-bind"))
	viper.BindPFlag("general.log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("general.skip_crc_check", rootCmd.PersistentFlags().Lookup("skip-crc-check"))
	viper.BindPFlag("backend.mqtt.server", rootCmd.PersistentFlags().Lookup("mqtt-server"))
	viper.BindPFlag("backend.mqtt.username", rootCmd.PersistentFlags().Lookup("mqtt-username"))
	viper.BindPFlag("backend.mqtt.password", rootCmd.PersistentFlags().Lookup("mqtt-password"))
	viper.BindPFlag("backend.mqtt.ca_cert", rootCmd.PersistentFlags().Lookup("mqtt-ca-cert"))
	viper.BindPFlag("backend.mqtt.tls_cert", rootCmd.PersistentFlags().Lookup("mqtt-tls-cert"))
	viper.BindPFlag("backend.mqtt.tls_key", rootCmd.PersistentFlags().Lookup("mqtt-tls-key"))

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)

}

func run(cmd *cobra.Command, args []string) error {
	log.SetLevel(log.Level(uint8(C.General.LogLevel)))

	log.WithFields(log.Fields{
		"version": version,
		"docs":    "https://docs.loraserver.io/lora-gateway-bridge/",
	}).Info("starting LoRa Gateway Bridge")

	var pubsub *mqttpubsub.Backend
	for {
		var err error
		pubsub, err = mqttpubsub.NewBackend(C.Backend.MQTT.Server, C.Backend.MQTT.Username, C.Backend.MQTT.Password, C.Backend.MQTT.CACert, C.Backend.MQTT.TLSCert, C.Backend.MQTT.TLSKey)
		if err == nil {
			break
		}

		log.Errorf("could not setup mqtt backend, retry in 2 seconds: %s", err)
		time.Sleep(2 * time.Second)
	}
	defer pubsub.Close()

	onNew := func(mac lorawan.EUI64) error {
		return pubsub.SubscribeGatewayTX(mac)
	}

	onDelete := func(mac lorawan.EUI64) error {
		return pubsub.UnSubscribeGatewayTX(mac)
	}

	gw, err := gateway.NewBackend(C.General.UDPBind, onNew, onDelete, C.General.SkipCRCCheck)
	if err != nil {
		log.Fatalf("could not setup gateway backend: %s", err)
	}
	defer gw.Close()

	go func() {
		for rxPacket := range gw.RXPacketChan() {
			if err := pubsub.PublishGatewayRX(rxPacket.RXInfo.MAC, rxPacket); err != nil {
				log.Errorf("could not publish RXPacket: %s", err)
			}
		}
	}()

	go func() {
		for stats := range gw.StatsChan() {
			if err := pubsub.PublishGatewayStats(stats.MAC, stats); err != nil {
				log.Errorf("could not publish GatewayStatsPacket: %s", err)
			}
		}
	}()

	go func() {
		for txPacket := range pubsub.TXPacketChan() {
			if err := gw.Send(txPacket); err != nil {
				log.Errorf("could not send TXPacket: %s", err)
			}
		}
	}()

	go func() {
		for txAck := range gw.TXAckChan() {
			if err := pubsub.PublishGatewayTXAck(txAck.MAC, txAck); err != nil {
				log.Errorf("could not publish TXAck: %s", err)
			}
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	log.Warning("shutting down server")
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "lora-gateway-bridge",
	Short: "abstracts the packet_forwarder protocol into JSON over MQTT",
	Long: `LoRa Gateway Bridge abstracts the packet_forwarder protocol into JSON over MQTT
	> documentation & support: https://docs.loraserver.io/lora-gateway-bridge
	> source & copyright information: https://github.com/brocaar/lora-gateway-bridge`,
	RunE: run,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the LoRa Gateway Bridge version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the LoRa Gateway configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := template.Must(template.New("config").Parse(configTemplate))
		err := t.Execute(os.Stdout, C)
		if err != nil {
			return errors.Wrap(err, "execute config template error")
		}
		return nil
	},
}

func initConfig() {
	if cfgFile != "" {
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}
		if err := viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}
	} else {
		viper.SetConfigName("lora-gateway-bridge")
		viper.AddConfigPath("/etc/lora-gateway-bridge/")
		viper.AddConfigPath("$HOME/.config/lora-gateway-bridge")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			log.WithError(err).Warning("falling back on environment variable configuration, update your configuration!")
		}
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.WithError(err).Fatal("unmarshal config error")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
