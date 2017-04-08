package mqtt

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/surgemq/message"
	"github.com/suyashkumar/conduit/server/secrets"
	"github.com/suyashkumar/surgemq/service"
)

var phaoClient MQTT.Client
var handlerMap = make(map[string]func(string, string))

func Register(name string, f func(string, string)) {
	handlerMap[name] = f
}

func DeRegister(name string) error {
	_, ok := handlerMap[name]
	// Mostly just a sanity check for callers:
	if !ok {
		return errors.New("Name never registered")
	}
	delete(handlerMap, name)
	return nil
}

var LOGGING = (os.ExpandEnv("LOGGING") != "")

func onPublish(client MQTT.Client, msg MQTT.Message) {
	if LOGGING {
		fmt.Println("Topic:", string(msg.Topic()), "Payload:", string(msg.Payload()))
	}

	if val, ok := handlerMap[string(msg.Topic())]; ok {
		val(string(msg.Topic()), string(msg.Payload()))
	}
	// Look to see if the published message was a streaming data message
	// If so, persist the contents to an appropiate db
	var validDataStream = regexp.MustCompile(`^[^/]*/stream/.*`)
	if validDataStream.MatchString(string(msg.Topic())) {
		go PersistMessage(string(msg.Payload()), string(msg.Topic()))
	}
}
func createServerClient() MQTT.Client {
	service.AllowedMap[secrets.SubSecret] = 1

	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID(secrets.SubSecret)

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("#", 0, onPublish); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	return c
}
func stayAlive(c *service.Client, KeepAlive int) {
	for _ = range time.Tick(time.Duration(KeepAlive) * time.Second) {
		c.Ping(func(msg, ack message.Message, err error) error {
			return nil
		})
	}
}

func sendMessage(c MQTT.Client, device string, payload string) {
	token := c.Publish(device, 0, false, payload)
	token.Wait()
}

func SendMessage(device string, payload string) {
	sendMessage(phaoClient, device, payload)
}

func RunServer() {
	fmt.Println("Starting up MQTT machinery...")
	svr := &service.Server{}
	go svr.ListenAndServe("tcp://:1883")
	time.Sleep(200 * time.Millisecond)
	phaoClient = createServerClient()
	fmt.Println("Started and listening")
}
