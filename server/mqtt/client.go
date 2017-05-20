package mqtt

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/suyashkumar/conduit/server/models"
	mgo "gopkg.in/mgo.v2"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/surgemq/message"
	"github.com/suyashkumar/conduit/server/secrets"
	"github.com/suyashkumar/surgemq/service"
)

var globalClient Client // The global instance of the internal mqtt globalClient

type Client interface {
	// SendMessage sends a message to a given device stream
	SendMessage(string, string)
	// Register registers a callback to fire on messages recieved on a given stream
	Register(string, func(string, string))
	// DeRegister removes a registered callback for a given device stream
	DeRegister(string) error
}

type client struct {
	MQTT        MQTT.Client
	CallbackMap map[string]func(string, string)
}

// GetClient returns pointer to current global instance of GetClient
func GetClient() Client {
	return globalClient
}

// Register registers a callback to fire on messages recieved on a given stream
func (c *client) Register(name string, f func(string, string)) {
	c.CallbackMap[name] = f
}

// DeRegister removes a registered callback for a given device stream
func (c *client) DeRegister(name string) error {
	_, ok := c.CallbackMap[name]
	// Mostly just a sanity check for callers:
	if !ok {
		return errors.New("Name never registered")
	}
	delete(c.CallbackMap, name)
	return nil
}

// SendMessage sends a message to a given device stream
func (c *client) SendMessage(device string, payload string) {
	token := c.MQTT.Publish(device, 0, false, payload)
	token.Wait()
}

var LOGGING = (os.ExpandEnv("LOGGING") != "")

func getPublishCallback(m map[string]func(string, string)) func(MQTT.Client, MQTT.Message) {
	return func(client MQTT.Client, msg MQTT.Message) {
		if LOGGING {
			fmt.Println("Topic:", string(msg.Topic()), "Payload:", string(msg.Payload()))
		}

		if val, ok := m[string(msg.Topic())]; ok {
			val(string(msg.Topic()), string(msg.Payload()))
		}
		// Look to see if the published message was a streaming data message
		// If so, persist the contents to an appropiate db
		var validDataStream = regexp.MustCompile(`^[^/]*/stream/.*`)
		if validDataStream.MatchString(string(msg.Topic())) {
			go persistMessage(string(msg.Payload()), string(msg.Topic()))
		}
	}
}

func persistMessage(message string, topic string) {
	session, err := mgo.Dial(secrets.DB_DIAL_URL)
	if err != nil {
		fmt.Println("ERROR Connecting to the database.", err)
	}
	defer session.Close()

	c := session.DB("homeauto").C("streammessages")
	err = c.Insert(&models.StreamMessage{
		Data:      message,
		Timestamp: time.Now(),
		Topic:     topic,
	})
	if err != nil {
		fmt.Println("ERROR inserting StreamMessage to database.", err)
	}
}

func createServerClient() Client {
	service.AllowedMap[secrets.SubSecret] = 1

	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID(secrets.SubSecret)

	c := client{
		MQTT:        MQTT.NewClient(opts),
		CallbackMap: make(map[string]func(string, string)),
	}

	if token := c.MQTT.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.MQTT.Subscribe("#", 0, getPublishCallback(c.CallbackMap)); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	return &c
}

func stayAlive(c *service.Client, keepAlive int) {
	for _ = range time.Tick(time.Duration(keepAlive) * time.Second) {
		c.Ping(func(msg, ack message.Message, err error) error {
			return nil
		})
	}
}

func sendMessage(c MQTT.Client, device string, payload string) {
	token := c.Publish(device, 0, false, payload)
	token.Wait()
}

func RunServer() {
	fmt.Println("Starting up MQTT machinery...")
	svr := &service.Server{}
	go svr.ListenAndServe("tcp://:1883")
	time.Sleep(200 * time.Millisecond)
	globalClient = createServerClient()
	fmt.Println("Started and listening")
}
