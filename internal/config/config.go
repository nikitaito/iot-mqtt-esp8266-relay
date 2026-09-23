package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

	ids, err := parseAllowedIDs("ALLOWED_TELEGRAM_USER_IDS")
	if err != nil {
		return nil, err
	}
	cfg.AllowedTelegramUserIDs = ids

	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseAllowedIDs(raw string) ([]int64, error) {
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid Telegram user ID %q: %w", p, err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (c *Config) validate() error {
	var missing []string

	if c.TelegramBotToken == "" {
		missing = append(missing, "TELEGRAM_BOT_TOKEN")
	}
	if c.MQTTBrokerURL == "" {
		missing = append(missing, "MQTT_BROKER_URL")
	}
	if c.MQTTTopic == "" {
		missing = append(missing, "MQTT_TOPIC")
	}
	if len(c.AllowedTelegramUserIDs) == 0 {
		missing = append(missing, "ALLOWED_TELEGRAM_USER_IDS")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}
	return nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
