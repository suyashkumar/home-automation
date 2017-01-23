package service

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/suyashkumar/conduit/server/mqtt"
	"github.com/suyashkumar/conduit/server/routes"
	"net/http"
	"os"
)

type ConduitService interface {
	Run()
}

type ConduitServiceImpl struct {
	Router *httprouter.Router
	IsDev  bool
}

func (c *ConduitServiceImpl) Run() {
	c.attachRoutes()
	mqtt.RunServer()
	c.startWebServer()
}

func (c *ConduitServiceImpl) attachRoutes() {
	c.Router.GET("/api/send/:deviceName/:funcName", routes.AuthMiddlewareGenerator(routes.Send))
	c.Router.GET("/api/streams/:deviceName/:streamName", routes.AuthMiddlewareGenerator(routes.GetStreamedMessages))
	c.Router.POST("/api/auth", routes.Auth)
	c.Router.POST("/api/register", routes.New)
	c.Router.GET("/api/me", routes.AuthMiddlewareGenerator(routes.GetUser))
	c.Router.GET("/", routes.Hello)
	c.Router.OPTIONS("/api/*sendPath", routes.Headers)
	c.Router.ServeFiles("/static/*filepath", http.Dir("public/static"))
}

func (c *ConduitServiceImpl) startWebServer() {
	fmt.Printf("Web server to listen on port :%s", os.Getenv("PORT"))
	if c.IsDev {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), c.Router)
		panic(err)
	} else {
		go http.ListenAndServeTLS(":443", os.Getenv("CERT"), os.Getenv("PRIV_KEY"), c.Router)
		err := http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(routes.RedirectToHttps))
		panic(err)
	}

}

func NewConduitService() *ConduitServiceImpl {
	return &ConduitServiceImpl{
		Router: httprouter.New(),
		IsDev:  os.Getenv("DEV") == "TRUE",
	}
}
