package device

import (
	"net/http"

	"fmt"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/suyashkumar/golang-socketio"
	"github.com/suyashkumar/golang-socketio/transport"
)

type Handler interface {
	Call(deviceName, deviceID, functionName string, wait bool) chan string
	On(deviceName, deviceID, eventName string, callback func(deviceName, eventName, body string))
	GetHTTPHandler() http.Handler
}

type handler struct {
	server *gosocketio.Server
}

var globalDeviceHandler *handler

func NewHandler() Handler {
	s := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	// Attach socket event handlers
	s.On("hello", onHello)
	s.On(gosocketio.OnConnection, onConnection)
	s.On("api_key", onAPIKeyReceive)

	globalDeviceHandler = &handler{
		server: s,
	}
	return globalDeviceHandler
}

func getRoomName(deviceName, deviceID string) string {
	return fmt.Sprintf("%s_%s", deviceID, deviceName)
}

func (h *handler) Call(deviceName, deviceID, functionName string, wait bool) chan string {
	reqUUID := uuid.NewV4().String()
	message := fmt.Sprintf("%s,%s", functionName, reqUUID)
	c := make(chan string)
	logrus.WithField("request_uuid", reqUUID).Info("Setting up event listener")
	h.server.On(reqUUID, func(ch *gosocketio.Channel, msg string) string {
		logrus.WithField("request_uuid", reqUUID).Info("Response returned")
		c <- msg
		return "OK"
	})
	h.server.BroadcastTo(getRoomName(deviceName, deviceID), "server_directives", message)
	return c
}

func (h *handler) On(deviceName, deviceID, eventName string, callback func(deviceName, eventName, body string)) {

}

func (h *handler) GetHTTPHandler() http.Handler {
	return h.server
}
