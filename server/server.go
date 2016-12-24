package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/suyashkumar/conduit/server/mqtt"
	"github.com/suyashkumar/conduit/server/routes"
	"net/http"
	"os"
)

func main() {
	router := httprouter.New()
	router.GET("/api/send/:deviceName/:funcName", routes.AuthMiddlewareGenerator(routes.Send))
	router.GET("/api/streams/:deviceName/:streamName", routes.AuthMiddlewareGenerator(routes.GetStreamedMessages))
	router.GET("/api/users", routes.ListUsers)
	router.POST("/api/auth", routes.Auth)
	router.POST("/api/register", routes.New)
	router.GET("/api/auth/test", routes.AuthMiddlewareGenerator(routes.Test))
	router.GET("/", routes.Hello)
	router.OPTIONS("/api/*sendPath", routes.Headers)

	mqtt.RunServer()
	fmt.Printf("Web server to listen on port :%s", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
	panic(err)
}
