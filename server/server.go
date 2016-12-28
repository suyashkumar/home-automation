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
	router.POST("/api/auth", routes.Auth)
	router.POST("/api/register", routes.New)
	router.GET("/api/me", routes.AuthMiddlewareGenerator(routes.GetUser))
	router.GET("/", routes.Hello)
	router.OPTIONS("/api/*sendPath", routes.Headers)
	router.ServeFiles("/static/*filepath", http.Dir("public/static"))

	mqtt.RunServer()
	fmt.Printf("Web server to listen on port :%s", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
	panic(err)
}
