package device

import (
	"sync"

	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/mqtt"
)

type State string

const (
	StateOn  State = "ON"
	StateOff State = "OFF"
)

type Relay struct {
	mqttClient *mqtt.Client
	mu         sync.Mutex
	state      State
}

func New(mqttclient *mqtt.Client) *Relay {
	return &Relay{
		mqttClient: mqttclient,
	}
}

func (rel *Relay) TurnON() error {
	err := rel.mqttClient.Publish("Relay_ON")
	if err != nil {
		return err
	}

	rel.mu.Lock()
	rel.state = StateOn
	rel.mu.Unlock()

	return nil
}

func (rel *Relay) TurnOFF() error {
	err := rel.mqttClient.Publish("Relay_OFF")
	return err
}
