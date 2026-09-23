package telegram

import (
	"github.com/go-telegram/bot"
	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/config"
	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/device"
)

type Bot struct {
	api   *bot.Bot
	cfg   *config.Config
	relay *device.Relay
}
