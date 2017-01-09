# conduit

Conduit allows you to quickly build cloud-connected hardware that you can control and communicate with from anywhere in the world. Conduit provides a RESTful API that allows you to easily call arbitrary functions (e.g. `int lightsOn()`) or to recieve/store data from the low-cost ESP8266 WiFi microcontroller. 

With Conduit you can:

- Dispatch ESP8266 firmware function calls on the target device via a RESTful API on the central conduit server (`GET https://conduit.suyash.io/api/send/:deviceName/:functionName`)
- Publish arbitrary data from the ESP8266 device to the conduit server (`conduit.publishData("hello", "testStream")` in the firmware) 
- Retreive previously published data via the simple RESTful API

all with minimal boilerplate and minimal setup :).

### Getting Started
Controlling an LED from the Cloud takes less than 5 minutes with Conduit. 

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
6. You should be set! You can now go to the conduit interact view (https://conduit.suyash.io/#/interact) and type in your device name (that you chose in step 4) and `ledToggle` as the function and hit "Go!" to see your LED on your device turn on :). 
7. There's a lot more to explore--you can publish persisted data to conduit (to be retrieved later via API) and build your own applications around conduit using the secure JSON web token based API.

### Sample Project
[smart-lights](https://github.com/suyashkumar/smart-lights) is a sample project that uses this library to switch lights from the cloud. 
![](https://github.com/suyashkumar/smart-lights/blob/master/img/lightswitch.gif)

### Example (In progress)
The basic functionality of this library is straightforward. Start with the provided platformio firmware template and just do the following: 
  
  1. Write a C function that returns an integer in your Arduino Code:
  
  ```C
  int ledOn(){
    digitalWrite(LED, HIGH);
    homeAuto.publish("LED is now on");
  }
  ```
  2. In your `setup()` function, register your function with the service: 
  
  ```C
  homeAuto.addHandler("ledON", &ledOn);
  ```

  3. [THIS WILL NOT WORK WITHOUT AN ACCOUNT or unless you're running the server locally] Go to `http://home.suyash.io/send/:deviceName/:functionString` and your function will run and any publish messages will be returned in the request response. Note, you must have a valid json web token in a `x-access-token` header. The :deviceName is set when initializing the HomeAuto object (line 27).
