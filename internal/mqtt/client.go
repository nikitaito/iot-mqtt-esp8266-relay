package mqtt

import (
	paho "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	pahoClient paho.Client
	topic      string
}
