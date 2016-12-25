/*
server.ino
Example for my library that handles ESP8266 communication with a server (even on private
networks). Consumers of this library can simply write functions and have them
be fired whenver the server fires a given event directed at this device. There is
a 1-1 mapping of event to function. For example the "led" event may fire the
ledToggle function on the device. The communication needed to get that event to the
device and decide what funciton to all is abstracted away entirely by this library.

@author: Suyash Kumar <suyashkumar2003@gmail.com>
*/
#include <Arduino.h>
#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <wifi_info.h> // comment this out and fill in the below two lines
#include <PubSubClient.h>
#include <Conduit.h>

#define LED 4

// Fill out the below Github peeps:
//const char* ssid = "mywifi";
//const char* password = "";

WiFiClient client;
PubSubClient pClient(client);
//HomeAuto conduit("suyash1", "home.suyash.io"); // or "suyash", "home.suyash.io"
//Conduit conduit("suyash1", "10.0.0.225"); // or "suyash", "home.suyash.io"
Conduit conduit("suyash", "192.168.1.144", "zHqHR0nSBTrIzaAY3JCY510Z"); // or "suyash", "home.suyash.io"
//Conduit conduit("suyash", "conduit.suyash.io", "a"); // or "suyash", "home.suyash.io"
int ledStatus = 0;

void startWIFI(){
  WiFi.begin(ssid, password);
  Serial.println("");

  // Wait for connection
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("");
  Serial.print("Connected to ");
  Serial.println(ssid);
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

}

int ledToggle(){
  digitalWrite(LED, (ledStatus) ? LOW : HIGH);
  ledStatus = (ledStatus) ? 0 : 1;
  Serial.println("Toggled");
  conduit.publishMessage((ledStatus) ? "LED ON" : "LED OFF");
}

int publishMessage(){
    conduit.publishMessage("hey there");
}

int publishSomeData(){
	conduit.publishData("10", "testing");
	conduit.publishMessage("Done");
}

void setup(void){
  Serial.begin(115200); // Start serial
  pinMode(LED, OUTPUT); // Set LED pin to output

  startWIFI(); // Config/start wifi

  // HomeAuto bindings
  conduit.addHandler("ledToggle", &ledToggle);
  conduit.addHandler("hello", &publishMessage);
  conduit.addHandler("publishSomeData", &publishSomeData);
  conduit.setClient(pClient);

}

void loop(void){
  conduit.handle();
}
