package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/suyashkumar/auth"
	"github.com/suyashkumar/conduit/db"
	"github.com/suyashkumar/conduit/device"
	"github.com/suyashkumar/conduit/handlers"
)

type deviceAPIHandler func(
	w http.ResponseWriter,
	r *http.Request,
	p httprouter.Params,
	d device.Handler,
	db db.Handler,
	a auth.Authenticator,
)

// injectMiddleware is middleware that injects the device.Handler into the RESTful API route handler functions
func injectMiddleware(h deviceAPIHandler, d device.Handler, db db.Handler, a auth.Authenticator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		h(w, r, p, d, db, a)
	}
}

// Build returns a new Router configured to serve all application routes
func Build(d device.Handler, db db.Handler, a auth.Authenticator) *httprouter.Router {
	r := httprouter.New()

	// Static serving routes:
	r.GET("/", handlers.Index)

	// RESTful API Routes:
	r.POST("/api/register", injectMiddleware(handlers.Register, d, db, a))
	r.POST("/api/login", injectMiddleware(handlers.Login, d, db, a))
	r.POST("/api/call", injectMiddleware(handlers.Call, d, db, a))

	// Configure static file serving from /static
	r.ServeFiles("/static/*filepath", http.Dir("public/static"))

	// Configure device handler socket routing
	r.Handler("GET", "/socket.io/", d.GetHTTPHandler())

	return r
}
