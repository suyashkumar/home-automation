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
// #include <wifi_info.h> // You can include secret wifi info in a seperate file 
#include <Conduit.h>

#define LED D0

// Fill out the below Github folks:
//const char* ssid = "mywifi";
//const char* password = "";
const char* deviceName = "suyasha";
const char* apiKey = "nTyiMTdyjS8UD1qA5hw4hSP1";
const char* serverUrl = "conduit.suyash.io";

Conduit conduit(deviceName, serverUrl, apiKey); // init Conduit 
int ledStatus = 0;

// Toggles an LED attached on the LED pin!
int ledToggle(){
  digitalWrite(LED, (ledStatus) ? LOW : HIGH);
  ledStatus = (ledStatus) ? 0 : 1;
  Serial.println("Toggled");
  conduit.publishMessage((ledStatus) ? "LED ON" : "LED OFF");
}

// Publishes a message response to the server 
// when this function is called
int publishMessage(){
    conduit.publishMessage("hey there");
}

// When this function is called
// sends data to the "testing" datastream
// to be persisted in a database on the server
// sends a "Done" response when done
int publishSomeData(){
	conduit.publishData("10", "testing");
	conduit.publishMessage("Done");
}

void setup(void){
  Serial.begin(115200); // Start serial
  pinMode(LED, OUTPUT); // Set LED pin to output
  digitalWrite(LED, LOW);

  conduit.startWIFI(ssid, password); // Config/start wifi
  conduit.init();

  // Conduit bindings allow you to use the
  // function name to call the associated function
  // using the conduit API
  conduit.addHandler("ledToggle", &ledToggle);
  conduit.addHandler("hello", &publishMessage);
  conduit.addHandler("publishSomeData", &publishSomeData); 

}

void loop(void){
  conduit.handle();
}
