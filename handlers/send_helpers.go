package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/suyashkumar/conduit/server/entities"
)

func sendOK(w http.ResponseWriter) {
	sendJSON(w, entities.GenericResponse{Message: "OK"}, 200)
}

func sendJSON(w http.ResponseWriter, v interface{}, statusCode int) error {
	resBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, string(resBytes))
	return nil
}
