package device

import (
	"github.com/sirupsen/logrus"
	"github.com/suyashkumar/golang-socketio"
)

const OK_MSG = "OK"

func onHello(c *gosocketio.Channel) string {
	logrus.Println("Something successfully handled")
	c.Emit("hello", "Hello emit")
	return "OK"
}

func onConnection(c *gosocketio.Channel) {
	logrus.Printf("New Connection (SID: %s)", c.Id())
	c.Emit("id_message", c.Id())
}

func onAPIKeyReceive(c *gosocketio.Channel, msg string) string {
	logrus.Infof("Received an API key message from %s: %s", c.Id(), msg)
	//TODO: Validate msg, consider receiving as JSON based on firmware
	c.Join(msg)
	return OK_MSG
}
