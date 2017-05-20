package main

import (
	"os"

	"fmt"

	"github.com/suyashkumar/conduit/server/app"
	"github.com/suyashkumar/conduit/server/secrets"
)

func main() {
	// Create Conduit App configuration
	c := app.Config{
		IsDev:     os.Getenv("DEV") == "TRUE",
		Port:      os.Getenv("PORT"),
		CertKey:   os.Getenv("CERT"),
		PrivKey:   os.Getenv("PRIV_KEY"),
		DBDialURL: secrets.DB_DIAL_URL,
	}

	// Create new app
	app, err := app.New(c)
	if err != nil {
		fmt.Println("Issue creating New app")
		panic(err)
	}

	// Run app
	err = app.Run()
	if err != nil {
		fmt.Println("Issue running app")
		panic(err)
	}
}
