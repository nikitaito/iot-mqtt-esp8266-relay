package device

import "github.com/nikitaito/iot-mqtt-esp8266-relay/internal/mqtt"

type State string

const (
	StateOn  State = "ON"
	StateOff State = "OFF"
)

type Relay struct {
	mqttClient *mqtt.Client
}

func New(mqttclient *mqtt.Client) *Relay {
	return &Relay{
		mqttClient: mqttclient,
	}
}

func (rel *Relay) TurnON() error {
	err := rel.mqttClient.Publish("Relay_ON")
	return err
}

func (rel *Relay) TurnOFF() error {
	err := rel.mqttClient.Publish("Relay_OFF")
	return err
}
