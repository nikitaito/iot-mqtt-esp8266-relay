# iot-mqtt-esp8266-relay

Control a relay (e.g. a light) from anywhere in the world, over Telegram, using an ESP8266 and a public MQTT broker — no self-hosted server required.

## How it works

```
Telegram (you) --/on /off--> Go controller --MQTT publish--> Broker --MQTT subscribe--> ESP8266 --> Relay
```

- A **Go program** (`cmd/controller`) runs a Telegram bot. When an authorized user sends `/on` or `/off`, it publishes a message to an MQTT topic.
- An **ESP8266** runs firmware that stays connected to the same broker, subscribed to the same topic, and switches a relay pin HIGH/LOW when it receives `Relay_ON` / `Relay_OFF`.
- Both sides talk through a **public, free MQTT broker** — nothing needs to be self-hosted or port-forwarded.

## ⚠️ About the public broker

This project is set up to use a free public broker (`broker.freemqtt.com`) with shared credentials (`freemqtt` / `public`). That means:

- **The username and password do not protect you** — anyone can connect to this broker with the same credentials.
- **The MQTT topic is your only real protection.** Whoever knows your topic name can publish `Relay_ON`/`Relay_OFF` to it and control your relay, or subscribe to it and see every command you send.

**You must use a long, random, hard-to-guess topic** — not something like `home/relay/set`. Example:

```
x7k2m9p4q8w1/relay1/set
```

Generate your own random string and use it consistently everywhere below. If you ever want real security (private credentials, access control), move to a broker like HiveMQ Cloud, which has a free tier with per-account auth.

## Repository layout

```
cmd/controller/         Go program entry point (main.go)
internal/config/        Loads and validates environment variables
internal/mqtt/          MQTT client wrapper (connect, publish, disconnect)
internal/device/        Relay abstraction (TurnON / TurnOFF)
internal/telegram/      Telegram bot command handlers
firmware/esp8266_relay/ Arduino sketch (.ino) that runs on the ESP8266
```

## Part 1 — Flashing the ESP8266

### Requirements

1. **Arduino IDE** (or Arduino CLI / PlatformIO, if you prefer).
2. **ESP8266 board support** installed in the Arduino IDE:
   - File → Preferences → Additional Board Manager URLs → add:
     `http://arduino.esp8266.com/stable/package_esp8266com_index.json`
   - Tools → Board → Boards Manager → search "esp8266" → install.
3. **PubSubClient library** (by Nick O'Leary):
   - Sketch → Include Library → Manage Libraries → search "PubSubClient" → install.

### Configure the sketch

Open `firmware/esp8266_relay/esp8266_relay.ino` and edit the top of the file:

```cpp
const char* ssid = "WIFI_SSID";
const char* password = "WIFI_PASSWORD";

const char* topic = "x7k2m9p4q8w1/relay1/set"; // your random topic — must match MQTT_TOPIC in .env
```

Everything else (broker host, port, broker username/password, payload strings) is already set to match this project's Go side — leave it as is unless you change the broker on both sides.

If you're wiring more than one ESP8266 for this project, give each one a unique `mqtt_client_id` (they can't share one — the broker disconnects duplicates).

### Upload

1. Connect the ESP8266 over USB.
2. Tools → Board → select your exact ESP8266 board (e.g. "NodeMCU 1.0").
3. Tools → Port → select the correct serial port.
4. Click Upload.
5. Open the Serial Monitor at 115200 baud to confirm it connects to WiFi and to the MQTT broker, and see the subscribed topic printed.

Wire the relay's control pin to **D4** (as set by `Relay_Pin` in the sketch), and the relay module's power/ground appropriately for your relay board.

## Part 2 — Running the Go controller

### Requirements

- Go 1.26+ (see `go.mod`)
- A Telegram bot token from [@BotFather](https://t.me/BotFather)
- Your Telegram numeric user ID (get it from [@userinfobot](https://t.me/userinfobot) — send it `/start`)

### Install dependencies

```bash
go mod tidy
```

### Configure

Copy `.env.example` to `.env` and fill in real values:

```dotenv
TELEGRAM_BOT_TOKEN=123456789:AA...your-bot-token...
MQTT_BROKER_URL=tls://broker.freemqtt.com:8883
MQTT_USERNAME=freemqtt
MQTT_PASSWORD=public
MQTT_CLIENT_ID=telegram-relay-controller
MQTT_TOPIC=x7k2m9p4q8w1/relay1/set
ALLOWED_TELEGRAM_USER_IDS=123456789
```

- `MQTT_TOPIC` **must exactly match** the `topic` value in the `.ino` firmware file.
- `MQTT_CLIENT_ID` **must be different** from the ESP8266's `mqtt_client_id` — both connect to the broker at the same time, and duplicate client IDs get disconnected.
- `ALLOWED_TELEGRAM_USER_IDS` is comma-separated if you want more than one authorized user. Anyone not on this list is refused when they try `/on` or `/off`.

`.env` is loaded automatically at startup via `godotenv` — you don't need to `export` anything manually.

### Run

```bash
go run ./cmd/controller
```

Expected output:

```
mqtt: connected to broker
controller: starting, press Ctrl+C to stop
```

Stop it any time with `Ctrl+C` — it disconnects from the broker cleanly before exiting.

## Part 3 — Using the bot

In Telegram, open a chat with your bot and send:

| Command | Effect |
|---|---|
| `/start` | Shows a short help message |
| `/on`    | Publishes `Relay_ON` — turns the relay on |
| `/off`   | Publishes `Relay_OFF` — turns the relay off |

If your Telegram user ID isn't in `ALLOWED_TELEGRAM_USER_IDS`, the bot replies that you're not authorized and ignores the command.

## Known limitations

- **No state feedback.** The controller only knows what it *sent*, not the relay's actual current state. If the ESP8266 restarts or loses power, the controller has no way to know. A `/status` command and a retained MQTT state topic from the ESP8266 side would close this gap — not implemented yet.
- **Shared broker credentials.** See the security note above — topic secrecy is the only real protection right now.
- **Single relay.** The project currently assumes one relay on one topic. Multiple relays would need per-device topics and matching handler logic.

## Troubleshooting

- **`missing required config: ...` on startup** — one of the required `.env` values is empty, or `.env` isn't being found (make sure you're running `go run ./cmd/controller` from the project root, where `.env` lives).
- **Bot doesn't respond at all** — check the bot token is correct and that you've sent `/start` to the bot at least once (Telegram won't let bots message users who haven't initiated contact).
- **Bot replies but relay doesn't move** — check the ESP8266's Serial Monitor: is it connected to the broker? Is it subscribed to the exact same topic as `MQTT_TOPIC`? Payload strings (`Relay_ON`/`Relay_OFF`) are case-sensitive and must match exactly on both sides.