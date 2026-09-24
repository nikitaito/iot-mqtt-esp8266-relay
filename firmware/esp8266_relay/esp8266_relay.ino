#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include <cstring>

const char* ssid = "WIFI_SSID";
const char* password = "WIFI_PASSWORD";

const char* mqtt_server = "";
const int mqtt_port = 1883; 

const char* mqtt_user = "";
const char* mqtt_pass = "";

const char* mqtt_client_id = "";

WiFiClient espClient;
PubSubClient client(espClient);

const char* topic = "TOPIC";

const char ON_PER[] = "Relay_ON";
const char OFF_PER[] = "Relay_OFF";

#define Relay_Pin D4


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


bool Equal(const char* per, byte* payload, unsigned int length) {

  unsigned int perLength = strlen(per);

  if (length != perLength) {
    return false;
  }

  for (unsigned int i = 0; i < perLength; i++) {
    if (payload[i] != per[i]) {
      return false;
    }
  }

  return true;
}


void callback(char* topic, byte* payload, unsigned int length) {

  Serial.print("Message received: ");

  if (Equal(ON_PER, payload, length)) {

    digitalWrite(Relay_Pin, LOW);

    Serial.println("Relay ON !!!!");

  }
  else if (Equal(OFF_PER, payload, length)) {

    digitalWrite(Relay_Pin, HIGH);

    Serial.println("Relay OFF !!!!");

  }
  else {

    Serial.println("Undefined Message !!!!");

  }
}


void reconnect() {
  while (!client.connected()) {
    Serial.print("Connecting to MQTT...");

    if (client.connect(mqtt_client_id, mqtt_user, mqtt_pass)) {
      Serial.println("connected!");

      client.subscribe(topic);

      Serial.println("Subscribed to:");
      Serial.println(topic);

    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" retrying...");

      delay(5000);
    }
  }
}

void setup() {
  Serial.begin(115200);

  pinMode(Relay_Pin, OUTPUT);

  setup_wifi();

  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(callback);
}


void loop() {
  if (!client.connected()) {
    reconnect();
  }

  client.loop();
}