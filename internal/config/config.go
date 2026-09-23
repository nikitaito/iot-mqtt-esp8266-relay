package config

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
