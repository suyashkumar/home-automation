# conduit
:eyes:Conduit V2 was just realeased! :eyes:

<a href="https://travis-ci.org/suyashkumar/conduit" target="_blank"><img src="https://travis-ci.org/suyashkumar/conduit.svg?branch=master" /></a>

[Conduit featured on Hackaday!](http://hackaday.com/2017/01/17/servo-controlled-iot-light-switches/)

Conduit is an entirely open-source web service that allows you to quickly and easily call functions on your [ESP8266 IoT devices](https://www.amazon.com/HiLetgo-Version-NodeMCU-Internet-Development/dp/B010O1G1ES/ref=sr_1_3?ie=UTF8&qid=1483953570&sr=8-3&keywords=nodemcu+esp8266) from anywhere in the world (even if those devices are behind private networks). 

You can do all this simply by dropping in a few lines of code into your firmware and then issuing RESTful API requests to the conduit web service to call your firmware functions. Skip ahead to a [full minimal example](README.md/#bink-an-led-from-the-cloud-full-example) if you're ready to get started right away!

## Conduit Components
* [Conduit backend web service (here)](https://github.com/suyashkumar/conduit)
* [Conduit firmware library](https://github.com/suyashkumar/conduit-firmware-library)
* [Conduit frontend](https://github.com/suyashkumar/conduit-frontend)

## Conduit API
A central conduit API server is already deployed at https://api.conduit.suyash.io (should be used for all API routes) with a user-friendly front-end deployed at https://conduit.suyash.io. 

| Method | Route          | Sample Request                                                                                                                           | Notes                                                                                                                     |
|--------|----------------|------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------|
| POST   | /api/login     | ``` {   "email": "test@test.com",   "password": "test" }  ```                                                                            | Authenticate with Conduit, get issued a JWT                                                                               |
| POST   | /api/send      | ```{"token": "JWT token from login", "device_name": "myDeviceName", "function_name": "ledToggle", "wait_for_device_response": "true"}``` | Call a function (ledToggle) on one of your ESP8266 devices (named "myDeviceName" here)!                                   |
| POST   | /api/user_info | ```{"token": "JWT token from login"}```                                                                                                  | This returns information about your user account, including your account  secret which you must include in your firmware. |


## Sample Application using Conduit
[smart-lights](https://github.com/suyashkumar/smart-lights) is a sample project that uses this library to switch lights from the cloud. It currently uses V1 of this library and should be updated to V2 shortly. 

![](https://github.com/suyashkumar/smart-lights/blob/master/img/lightswitch.gif)

## Minimal Example
Below is a minimal example of firmware code needed to get started with Conduit to blink an LED. See the next section for a complete example. 
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

Conduit is currently in active development, so please feel free to contact me with comments/questions and submit pull requests!

## Bink an LED from the Cloud (full example).
Controlling an LED on the ESP8266 from the Cloud takes less than 5 minutes with Conduit. 

Please make sure you've installed the relevant drivers ([here](https://www.silabs.com/products/mcu/Pages/USBtoUARTBridgeVCPDrivers.aspx) if you're using the nodemcu ESP8266 chip linked above) and installed the [platformio](http://docs.platformio.org/en/latest/installation.html) build system (simply `brew install platformio` if you're on a mac).

1. Create a conduit account at https://conduit.suyash.io/#/login
2. Retreive your account secret from the Account view at https://conduit.suyash.io/#/account
3. Clone the conduit firmware repo and change into the `examples/basic_functionality` directory.

  ```sh
  git clone https://github.com/suyashkumar/conduit-firmware-library.git
  cd examples/basic_functionality
  ```
4. Open `src/main.ino`. Fill in the following lines (account secret comes from step 2):

  ```C
  // Fill out the below Github folks:
const char* ssid = "";
const char* password = "";
const char* device_name = "myDevice";
const char* server_url = "api.conduit.suyash.io";
const char* account_secret = "";
  ```
5. Build the project using platformio. You should [install platformio](http://docs.platformio.org/en/latest/installation.html#python-package-manager) (if you haven't already) to build this properly. Ensure you're in the root directory of the example (not `src`) and run:

  ```sh
  platformio run
  ```
  If your ESP8266 chip is connected via usb already, to build and upload the program run:
  ```sh
  platformio run --target upload
  ```
  NOTE: to properly upload to an ESP8266 chip, you must have installed the ESP8266 drivers on your system already.

6. You should be set! You can now go to the conduit interact view (https://conduit.suyash.io/#/interact) and type in your device name (that you chose in step 4) and `ledToggle` as the function and hit "Execute!" to see your LED on your device toggle! Note that because we're using the built-in LED the on/off statuses are reversed (LED is on when D0 is low), but with your own LED things should be normal!

## License 
Copyright (c) 2018 Suyash Kumar

See [conduit/LICENSE.txt](https://github.com/suyashkumar/conduit/blob/master/LICENSE.txt) for license text (CC Attribution-NonCommercial 3.0)
