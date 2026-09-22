#include <ESP8266WiFi.h>
#include <PubSubClient.h>

const char* ssid = "WIFI_SSID";
const char* password = "WIFI_PASSWORD";

const char* mqtt_server = "MQTT_SERVER_ADDRESS";

WiFiClient espClient;
PubSubClient client(espClient);

const char* topic = "TOPIC";
