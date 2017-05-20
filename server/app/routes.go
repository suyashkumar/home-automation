package app

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/suyashkumar/conduit/server/handlers"
)

func injectContext(h handlers.Handler, ctx *handlers.Context) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r, ps, ctx)
	}
}
func injectAuthMiddleware(next handlers.AuthHandler, c *handlers.Context) httprouter.Handle {

	middleware := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handlers.SetCorsHeaders(w)
		if candidateToken, ok := r.Header["X-Access-Token"]; ok {
			// Parse and validate token:
			token, err := jwt.ParseWithClaims(candidateToken[0], &handlers.HomeAutoClaims{}, func(token *jwt.Token) (interface{}, error) {
				return handlers.SecretKey, nil
			})

			if claims, ok := token.Claims.(*handlers.HomeAutoClaims); ok && token.Valid {
				next(w, r, ps, c, claims)
				return
			} else {
				handlers.SendErrorResponse(w, err.Error(), 401)
				fmt.Println("Error in Auth middleware")
				fmt.Println(err.Error())
				return
			}
		}
		// Either token wasn't valid or it wasn't provided
		handlers.SendErrorResponse(w, "No Token", 400)
		return
	}

	return middleware

}

func attachRoutes(r *httprouter.Router, ctx *handlers.Context) {
	r.GET("/api/send/:deviceName/:funcName", injectAuthMiddleware(handlers.Send, ctx))
	r.GET("/api/streams/:deviceName/:streamName", injectAuthMiddleware(handlers.GetStreamedMessages, ctx))
	r.POST("/api/auth", injectContext(handlers.Auth, ctx))
	r.POST("/api/register", injectContext(handlers.New, ctx))
	r.GET("/api/me", injectAuthMiddleware(handlers.GetUser, ctx))
	r.GET("/", handlers.Hello)
	r.OPTIONS("/api/*sendPath", handlers.Headers)
	r.ServeFiles("/static/*filepath", http.Dir("public/static"))
}
