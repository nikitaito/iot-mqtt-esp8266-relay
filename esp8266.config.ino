#include <ESP8266WiFi.h>
#include <PubSubClient.h>

const char* ssid = "WIFI_SSID";
const char* password = "WIFI_PASSWORD";

const char* mqtt_server = "MQTT_SERVER_ADDRESS";

WiFiClient espClient;
PubSubClient client(espClient);

const char* topic = "TOPIC";

void setup_wifi() {
  delay(10);

  WiFi.begin(ssid, password);

  Serial.print("Connecting to WiFi");

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println();
  Serial.println("WiFi connected!");
  Serial.println(WiFi.localIP());
}

