package mqtt

import (
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/config"
)

type Client struct {
	pahoClient paho.Client
	topic      string
}

func New(cfg *config.Config) (*Client, error) {
	opts := paho.NewClientOptions().
		AddBroker(cfg.MQTTBrokerURL).
		SetClientID(cfg.MQTTClientID).
		SetUsername(cfg.MQTTUsername).
		SetPassword(cfg.MQTTPassword).
		SetAutoReconnect(true).
		SetConnectTimeout(10 * time.Second).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second)

}
