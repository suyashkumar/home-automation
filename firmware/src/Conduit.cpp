/*
Conduit.cpp
A library that handles ESP8266 communication with a server (even on private
networks). Consumers of this library can simply write functions and have them
be fired whenver the server fires a given event directed at this device. There is
a 1-1 mapping of event to function. For example the "led" event may fire the
ledToggle function on the device. The communication needed to get that event to the
device and decide what funciton to all is abstracted away entirely by this library.

@author: Suyash Kumar <suyashkumar2003@gmail.com>
*/

#include "Conduit.h"

#define PREFIXED_NAME_CHARS   35
#define STATIC_STRING_CHARS   64
#define RECONNECT_DELAY_MS    5000
#define CONNECT_DELAY_MS      500
#define ROOT_NODE_NAME        "ROOT"

// These values are required for interop with Conduit (conduit.suyash.io)!
// DO NOT CHANGE THESE
#define DEVICE_PATH_SUFFIX    "/device"
#define STREAM_PATH_INFIX     "/stream/"
#define DEFAULT_TOPIC_NAME    "outTopic"
#define ON_CONNECT_MSG        "!connected"
#define SERVER_DEFAULT_PORT   1883  

#define STRCMP_SAME           0

typedef struct node {
  handler f;
  const char* name;
  node *next;
} node_t;

node_t* root;
node_t* current;
char prefixed_name[PREFIXED_NAME_CHARS];
WiFiClient client;
PubSubClient pClient(client);

Conduit::Conduit(const char* name, const char* server, const char* prefix){
  // Set name and server
  this->_name = name;
  this->_mqtt_server = server; 
  this->_prefix = prefix;

  // Compute _prefixed_name
  strcpy(prefixed_name, prefix);
  strcat(prefixed_name, name);
  const char* full_prefixed_name = prefixed_name;
  this->_prefixed_name = full_prefixed_name; 

  // Init linked list
  root = (node_t *)malloc(sizeof(node_t));
  root->next = nullptr;
  root->name = ROOT_NODE_NAME; // reserved name for root node (for now)
  current = root; 
}

void Conduit::addHandler(const char* name, handler f){
  node *newNode = (node_t*) malloc(sizeof(node_t));
  newNode->f = f;
  newNode->next = nullptr;
  newNode->name = name;
  current->next = newNode;
  current = newNode;
}

void Conduit::callHandler(const char* name){
  node_t* currentInSearch = root;
  while(true){
    if (strcmp(name, currentInSearch->name) == STRCMP_SAME){
      currentInSearch->f(); // Call function assoc with handler
      break;
    }
    if(currentInSearch->next == nullptr){
      break;
    }
    currentInSearch = currentInSearch->next;
  }
}

Conduit& Conduit::init(){
  this->_client = &pClient;
  pClient.setServer(this->_mqtt_server, SERVER_DEFAULT_PORT);
  Serial.println("Happened");
  pClient.setCallback([&](char* topic, byte* payload, unsigned int length){
    Serial.print("Message arrived [");
    Serial.print(topic);
    Serial.print("] ");
    char payloadStr[length];
    for (int i = 0; i < length; i++) {
      payloadStr[i] = (char)payload[i];
    }
    payloadStr[length] = 0; // Null terminate payloadStr
    removeSpace(payloadStr);
    Serial.println(payloadStr);
    this->callHandler(payloadStr); // Call function assoc with this handler
  });
  return *this;
}

void Conduit::handle(){
  if (!this->_client->connected()){
    this->reconnect();
  }
  this->_client->loop();
}

void Conduit::reconnect() {
  // Loop until we're reconnected
  while (!this->_client->connected()) {
    Serial.print("Attempting MQTT connection...");
    // Attempt to connect
	Serial.println(this->_prefixed_name);
    if (this->_client->connect(this->_prefixed_name)) {
      Serial.println("connected");
      // Once connected, publish an announcement...
      this->_client->publish(DEFAULT_TOPIC_NAME, ON_CONNECT_MSG);
      // Suscribe to topics:
      this->_client->subscribe(this->_prefixed_name); // suscribe to events meant for this device
    } else {
      Serial.print("failed, rc=");
      Serial.print(this->_client->state());
      Serial.println(" try again in 5 seconds"); 
      // Wait 5 seconds before retrying
      delay(RECONNECT_DELAY_MS);
    }
  }
}

void Conduit::publishMessage(const char* message){
  char str[STATIC_STRING_CHARS];
  strcpy(str, this->_prefixed_name);
  strcat(str, DEVICE_PATH_SUFFIX);
  const char* topicName = str;
  this->_client->publish(topicName, message);
}

void Conduit::publishData(const char* message, const char* dataStream) { 
	char topicBuffer[STATIC_STRING_CHARS];
	strcpy(topicBuffer, this->_prefixed_name);
	strcat(topicBuffer, STREAM_PATH_INFIX);
	strcat(topicBuffer, dataStream);
	const char* topicName = topicBuffer;
	this->_client->publish(topicName, message);
}

void Conduit::startWIFI(const char* ssid, const char* password){
  WiFi.begin(ssid, password);
  Serial.println("");

  // Wait for connection
  while (WiFi.status() != WL_CONNECTED) {
    delay(CONNECT_DELAY_MS);
    Serial.print(".");
  }

  Serial.println("");
  Serial.print("Connected to ");
  Serial.println(ssid);
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

}


void removeSpace(char* s) {
    for (char* s2 = s; *s2; ++s2) {
        if (*s2 != ' ')
            *s++ = *s2;
    }
    *s = 0;
}

