package routes

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.ServeFile(w, r, "public/index.html")
}

// Redirect to https://
func RedirectToHttps(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://conduit.suyash.io"+r.RequestURI, http.StatusMovedPermanently)
}
