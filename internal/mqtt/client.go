package mqtt

import (
	"fmt"
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
	opts.OnConnect = func(c paho.Client) {
		fmt.Println("mqtt : connected to brocker")
	}

	opts.OnConnectionLost = func(c paho.Client, err error) {
		fmt.Printf("mqtt: connection lost: %v\n", err)
	}

	pahoClient := paho.NewClient(opts)

	token := pahoClient.Connect()
	if !token.WaitTimeout(10 * time.Second) {
		return nil, fmt.Errorf("mqtt: connect timed out")
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("mqtt: connect failed: %w", err)
	}
}
