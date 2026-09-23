package telegram

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

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

func (b *Bot) isAllowed(ctx context.Context, update *models.Update) bool {
	if update.Message == nil || update.Message.From == nil {
		return false
	}

	if b.cfg.IsAllowed(update.Message.From.ID) {
		return true
	}

	b.reply(ctx, update.Message.Chat.ID, "You are not authorized to use this bot.")
	log.Printf("telegram: rejected command from unauthorized user id=%d", update.Message.From.ID)
	return false
}

func (b *Bot) reply(ctx context.Context, chatID int64, text string) {
	if _, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		log.Printf("telegram: failed to send message: %v", err)
	}
}
