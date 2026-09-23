package telegram

import (
	"context"

	"github.com/go-telegram/bot"

	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/config"
	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/device"
)

type Bot struct {
	api   *bot.Bot
	cfg   *config.Config
	relay *device.Relay
}

func New(cfg *config.Config, relay *device.Relay) (*Bot, error) {
	b := &Bot{
		cfg:   cfg,
		relay: relay,
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(b.handleUnknown),
	}

	api, err := bot.New(cfg.TelegramBotToken, opts...)
	if err != nil {
		return nil, err
	}
	b.api = api

	api.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.handleStart)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/on", bot.MatchTypeExact, b.handleOn)
	api.RegisterHandler(bot.HandlerTypeMessageText, "/off", bot.MatchTypeExact, b.handleOff)

	return b, nil
}

func (b *Bot) Start(ctx context.Context) {
	b.api.Start(ctx)
}
