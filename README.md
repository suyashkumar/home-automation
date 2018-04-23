# conduit
:eyes:Conduit V2 was just realeased! :eyes:

<a href="https://travis-ci.org/suyashkumar/conduit" target="_blank"><img src="https://travis-ci.org/suyashkumar/conduit.svg?branch=master" /></a>

[Conduit featured on Hackaday!](http://hackaday.com/2017/01/17/servo-controlled-iot-light-switches/)

Conduit is an open-source secure web service that allows you to quickly and easily call functionality on your ESP8266 IoT devices from anywhere in the world (even if those devices are behind private networks). 

You can do all this simply by dropping in a few lines of code into your firmware and then issuing RESTful API requests to the conduit web service to call your firmware functions. Suppose I wanted to toggle an LED on my ESP8266 in Florida from an iPhone app in North Carolina:

```C
#include <Arduino.h> 
#include <Conduit.h>

#define LED D0

const char* ssid = "ssid";
const char* password = "password";

Conduit conduit("myDeviceName", "api.conduit.suyash.io", "my-conduit-account-secret");
int ledStatus = 0;

int ledToggle(RequestParams *rq){
  digitalWrite(LED, (ledStatus) ? HIGH : LOW); // LED is on when LOW
  ledStatus = (ledStatus) ? 0 : 1;
  conduit.sendResponse(rq, (ledStatus) ? "ON":"OFF"); // send response to conduit
}


void setup(void){
  Serial.begin(115200); // Start serial
  pinMode(LED, OUTPUT); // Set LED pin to output
  digitalWrite(LED, HIGH);

  conduit.startWIFI(ssid, password); // Config/start wifi
  conduit.init();
  conduit.addHandler("ledToggle", &ledToggle); // register ledToggle as "ledToggle" with Conduit

}

void loop(void){
  conduit.handle();
}
```

and now you can call `ledToggle` on that device from anywhere in the world with:
  * POST https://api.conduit.suyash.io/api/send with the following JSON body:
    ```json
    {
      "token": "<your_jwt_token_here>",
      "device_name": "myDeviceName", 
      "function_name": "ledToggle",
      "wait_for_device_response": "true"
    }
    ```

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
Copyright (c) 2018 Suyash Kumar

See [conduit/LICENSE.txt](https://github.com/suyashkumar/conduit/blob/master/LICENSE.txt) for license text (CC Attribution-NonCommercial 3.0)
