package main

import (
	"github.com/suyashkumar/conduit/server/service"
)

func main() {
	conduitService := service.NewConduitService() // Init a new conduit web service
	conduitService.Run()                          // Run the conduit web service server
}
