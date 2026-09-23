package device

import "github.com/nikitaito/iot-mqtt-esp8266-relay/internal/mqtt"

type Relay struct {
	mqttClient *mqtt.Client
}
