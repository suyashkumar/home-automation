package routes

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/suyashkumar/conduit/server/mqtt"
	"net/http"
	"time"
)

type RpcResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
}

func PrefixedName(deviceName string, prefix string) string {
	return prefix + deviceName
}

func Send(w http.ResponseWriter, r *http.Request, ps httprouter.Params, hc *HomeAutoClaims) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		fmt.Println("OPT")
		return
	}

	mqtt.SendMessage(ps.ByName("deviceName"), ps.ByName("funcName"))
	c := make(chan string)
	end := make(chan string)

	mqtt.Register(ps.ByName("deviceName")+"/device", func(topic string, payload string) {
		res := &RpcResponse{
			Success: true,
			Data:    payload,
		}
		resBytes, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(resBytes))
		c <- "done"
	})
	time.AfterFunc(2*time.Second, func() {
		fmt.Fprintf(w, "ERROR")
		end <- "done"
	})
	select {
	case <-c:
		return
	case <-end:
		return
	}

}
