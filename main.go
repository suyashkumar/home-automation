package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/suyashkumar/auth"
	"github.com/suyashkumar/conduit/config"
	db2 "github.com/suyashkumar/conduit/db"
	"github.com/suyashkumar/conduit/device"
	"github.com/suyashkumar/conduit/log"
	"github.com/suyashkumar/conduit/routes"
)

func main() {
	log.Configure()

	d := device.NewHandler()
	db, err := db2.NewHandler(config.Get(config.DBConnString))
	if err != nil {
		logrus.WithError(err).WithField("DBConnString", config.Get(config.DBConnString)).Fatal("Could not connect to DB")
	}
	a, err := auth.NewAuthenticatorFromGORM(db.GetDB(), []byte(config.Get(config.SigningKey)))
	if err != nil {
		logrus.WithError(err).Fatal("Could not connect to or init database")
	}
	r := routes.Build(d, db, a)
	handler := cors.Default().Handler(r)

	p := fmt.Sprintf(":%s", config.Get(config.Port))

	if config.Get(config.UseSSL) == "false" {
		logrus.WithField("port", p).Info("Serving without SSL")
		err := http.ListenAndServe(p, handler)
		logrus.Fatal(err)
	} else {
		logrus.Info("Serving with SSL")
		err := http.ListenAndServeTLS(
			p,
			config.Get(config.CertKey),
			config.Get(config.PrivKey),
			handler,
		)
		// TODO: reroute http requests to https
		logrus.Fatal(err)
	}

}
