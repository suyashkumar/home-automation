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
2. Clone this repo and change into the conduit directory.

  ```sh
  git clone https://github.com/suyashkumar/conduit.git
  cd conduit
  ```
3. Navigate into the firmware directory and open `src/main.ino`

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
