package config

import (
	"errors"
	"os"
)

// Config holds everything the controller needs: Telegram bot credentials,
// MQTT broker connection info, and the list of Telegram user IDs allowed
// to control the relay.
type Config struct {
	TelegramBotToken string

	MQTTBrokerURL string
	MQTTClientID  string
	MQTTUsername  string
	MQTTPassword  string
	MQTTTopic     string

	AllowedTelegramUserIDs []int64
}

func Load() (*Config, error) {
	cfg := &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		MQTTBrokerURL:    os.Getenv("MQTT_BROKER_URL"),
		MQTTClientID:     getEnvOrDefault("MQTT_CLIENT_ID", "telegram-relay-controller"),
		MQTTUsername:     os.Getenv("MQTT_USERNAME"),
		MQTTPassword:     os.Getenv("MQTT_PASSWORD"),
		MQTTTopic:        os.Getenv("MQTT_TOPIC"),
	}
	return cfg, errors.New("")

}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
