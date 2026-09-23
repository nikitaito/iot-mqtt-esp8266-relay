package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/config"
	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/device"
	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/mqtt"
	"github.com/nikitaito/iot-mqtt-esp8266-relay/internal/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	mqttClient, err := mqtt.New(cfg)
	if err != nil {
		log.Fatalf("mqtt: %v", err)
	}
	defer mqttClient.Disconnect()

	relay := device.New(mqttClient)

	tgBot, err := telegram.New(cfg, relay)
	if err != nil {
		log.Fatalf("telegram: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("controller: starting, press Ctrl+C to stop")
	tgBot.Start(ctx)
	log.Println("controller: stopped")
}
