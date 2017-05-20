package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

type Context struct {
	DbSession *mgo.Session
}

type AuthHandler func(
	http.ResponseWriter,
	*http.Request,
	httprouter.Params,
	*Context,
	*HomeAutoClaims,
)

type Handler func(
	http.ResponseWriter,
	*http.Request,
	httprouter.Params,
	*Context,
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
