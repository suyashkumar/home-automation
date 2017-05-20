package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/suyashkumar/conduit/server/mqtt"
)

type RpcResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
}

//TODO: this should be in a name transforms package
func PrefixedName(deviceName string, prefix string) string {
	return prefix + deviceName
}

func Send(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *Context, hc *HomeAutoClaims) {

	prefixedName := PrefixedName(ps.ByName("deviceName"), hc.Prefix)

	mqtt.GetClient().SendMessage(prefixedName, ps.ByName("funcName"))

	c := make(chan string)

	mqtt.GetClient().Register(prefixedName+"/device", func(topic string, payload string) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Error in device response handler", r)
			}
		}()
		res := &RpcResponse{
			Success: true,
			Data:    payload,
		}
		resBytes, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(resBytes))
		c <- "done"

	})

	timeout := time.After(3 * time.Second)

	select {
	case <-c:
		// Got device response! Do nothing, move on
		fmt.Printf("Got device %s response", ps.ByName("deviceName"))
	case <-timeout:
		// Timed out!
		SendErrorResponse(w, "ERROR, no response from device", 504)
		fmt.Println("Timeout waiting for response from device")
	}

	// Cleanup:
	err := mqtt.GetClient().DeRegister(prefixedName + "/device")
	if err != nil {
		fmt.Println("Issues deregistering device message listener")
	}
}
