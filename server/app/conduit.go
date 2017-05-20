package app

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/suyashkumar/conduit/server/handlers"
	"github.com/suyashkumar/conduit/server/mqtt"
	"gopkg.in/mgo.v2"
)

type App interface {
	Run() error
}

type app struct {
	config  Config
	router  *httprouter.Router
	context *handlers.Context
}

type Config struct {
	IsDev     bool
	Port      string
	CertKey   string
	PrivKey   string
	DBDialURL string
}

func New(c Config) (*app, error) {
	// Initialize DB Session
	session, err := mgo.Dial(c.DBDialURL)
	if err != nil {
		return nil, err
	}

	// Initialize Context
	ctx := &handlers.Context{DbSession: session}

	// Initialize Router
	r := httprouter.New()
	attachRoutes(r, ctx)

	return &app{
		router: r,
		config: c,
	}, nil
}

func (a *app) Run() error {
	mqtt.RunServer()

	fmt.Printf("Web server to listen on port :%s", a.config.Port)
	if a.config.IsDev {
		return http.ListenAndServe(":"+a.config.Port, a.router)
	} else {
		go http.ListenAndServeTLS(":443", a.config.CertKey, a.config.PrivKey, a.router)
		return http.ListenAndServe(":"+a.config.Port, http.HandlerFunc(handlers.RedirectToHttps))
	}
}
