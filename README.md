# conduit
<a href="https://travis-ci.org/suyashkumar/conduit" target="_blank"><img src="https://travis-ci.org/suyashkumar/conduit.svg?branch=master" /></a>

[Conduit featured on Hackaday!](http://hackaday.com/2017/01/17/servo-controlled-iot-light-switches/)

Conduit allows you to quickly build cloud-connected IoT devices that you can communicate with and control from anywhere in the world. Conduit provides a RESTful API that allows you to remotely call functions (e.g. `lightsOn()`) on your ESP8266/Arduino device from the cloud, even if it's behind a LAN and doesn't have a public IP address. You can do all this simply by dropping in a few lines of code into your firmware:

```C
#include <Arduino.h>
#include <Conduit.h>

#define LED D0

Conduit conduit("my-device-name", "conduit.suyash.io", "api-key-here"); // init Conduit

// Turns on a LED
int ledOn() {
  digitalWrite(LED, HIGH);
  conduit.publishMessage("LED ON");
}

void setup(void) {
  pinMode(LED, OUTPUT); // Set LED pin to output
  digitalWrite(LED, LOW); // Start out with LED off

  conduit.startWIFI("ssid", "password"); // Config/start wifi
  conduit.init();
  conduit.addHandler("ledOn", &ledOn); // Registers ledOn function to be callable remotely
}

void loop(void) {
  conduit.handle();
}
```

...and just like that, you can now call `ledOn` on this device by making a `GET https://conduit.suyash.io/api/send/my-device-name/ledOn` from anywhere in the world! You will also have to supply a valid auth token by including an `x-access-token` header to ensure your call is authorized, and have an associated valid Conduit account. 

Conduit also provides a streamlined interface for recieving and making available arbitrary data produced from your devices in real time. 

Conduit is **entirely open source** (the firmware, backend web service, and frontend), allowing you to deploy your own instance of Conduit behind protected networks (like hospitals) or to audit the Conduit code. Conduit currently works with the [low-cost ESP8266 WiFi microcontroller](https://www.amazon.com/HiLetgo-Version-NodeMCU-Internet-Development/dp/B010O1G1ES/ref=sr_1_3?ie=UTF8&qid=1483953570&sr=8-3&keywords=nodemcu+esp8266) or Arduino like microcontroller, but there is no reason why it can't also work on other systems.

Conduit is currently in active development, so please feel free to contact me with comments/questions and submit pull requests!

### Bink an LED from the Cloud
Controlling an LED on the ESP8266 from the Cloud takes less than 5 minutes with Conduit. Please make sure you've installed the relevant drivers ([here](https://www.silabs.com/products/mcu/Pages/USBtoUARTBridgeVCPDrivers.aspx) if you're using the nodemcu ESP8266 chip linked above) and installed the [platformio](http://docs.platformio.org/en/latest/installation.html) build system (simply `brew install platformio` if you're on a mac).

1. Create a conduit account at https://conduit.suyash.io/#/login
2. Retreive your API key from the Account view at https://conduit.suyash.io/#/account
3. Clone this repo and change into the conduit directory.

  ```sh
  git clone https://github.com/suyashkumar/conduit.git
  cd conduit
  ```
4. Navigate into the firmware directory (`cd firmware`) and open `src/main.ino`. Fill in the following lines (API key comes from step 2):

  ```C
  // Fill out the below Github folks:
  const char* ssid = "mywifi";
  const char* password = "";
  const char* deviceName = "suyash";
  const char* apiKey = "api-key-here";
  ```
5. Build the project using platformio. You should [install platformio](http://docs.platformio.org/en/latest/installation.html#python-package-manager) (if you haven't already) to build this properly. Ensure you're in the firmware directory (`conduit/firmware`) and run:

  ```sh
  platformio run
  ```
  If your ESP8266 chip is connected via usb already, to build and upload the program run:
  ```sh
  platformio run --target upload
  ```
  NOTE: to properly upload to an ESP8266 chip, you must have installed the ESP8266 drivers on your system already.

6. You should be set! You can now go to the conduit interact view (https://conduit.suyash.io/#/interact) and type in your device name (that you chose in step 4) and `ledToggle` as the function and hit "Go!" to see your LED on your device toggle! Note that because we're using the built-in LED the on/off statuses are reversed (LED is on when D0 is low), but with your own LED things should be normal!
7. There's a lot more to explore--you can publish persisted data to conduit (to be retrieved later via API) and build your own applications around conduit using the secure JSON web token based API.

### Sample Project
[smart-lights](https://github.com/suyashkumar/smart-lights) is a sample project that uses this library to switch lights from the cloud. 
![](https://github.com/suyashkumar/smart-lights/blob/master/img/lightswitch.gif)

### License 
Copyright (c) 2017 Suyash Kumar

See [conduit/LICENSE.txt](https://github.com/suyashkumar/conduit/blob/master/LICENSE.txt) for license text (CC Attribution-NonCommercial 3.0)
