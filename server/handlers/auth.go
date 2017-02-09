package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/suyashkumar/conduit/server/models"
	"github.com/suyashkumar/conduit/server/secrets"
	"github.com/suyashkumar/conduit/server/util"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

var SecretKey = []byte(secrets.SECRET)

const JWT_TTL = 720      // In minutes
const PREFIX_LENGTH = 24 // Characters or bytes

type HomeAutoClaims struct {
	Email  string `json:"email"`
	Prefix string `json:"prefix"`
	jwt.StandardClaims
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type TokenResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

type UserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Email   string `json:"email"`
	Key     string `json:"key"`
}

type HandlerContext struct {
	DbSession *mgo.Session
}

type AuthHandler func(
	http.ResponseWriter,
	*http.Request,
	httprouter.Params,
	*HandlerContext,
	*HomeAutoClaims,
)

type ConduitHandler func(
	http.ResponseWriter,
	*http.Request,
	httprouter.Params,
	*HandlerContext,
)

func returnHash(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword)
}

func ListUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session, err := mgo.Dial(secrets.DB_DIAL_URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("homeauto").C("users")

	var results []models.User
	err = c.Find(bson.M{}).All(&results)
	if err != nil {
		panic(err)
	}
	resBytes, _ := json.Marshal(results)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(resBytes))
}

func New(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *HandlerContext) {
	SetCorsHeaders(w)
	u, err := decodeUserFromRequest(r)
	if err != nil {
		SendErrorResponse(w, err.Error(), 400)
		return
	}
	u.Prefix = util.GetRandString(PREFIX_LENGTH)
	u.Password = returnHash(u.Password)
	session := context.DbSession.New()
	defer session.Close()
	c := session.DB("homeauto").C("users")
	err = c.Insert(u)
	if err != nil {
		SendErrorResponse(w, err.Error(), 500)
	}
	fmt.Fprintf(w, "DONE")
}

func decodeUserFromRequest(r *http.Request) (models.User, error) {
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return u, err
	}
	// TODO: Add validation
	return u, nil
}

func SendErrorResponse(w http.ResponseWriter, errorString string, errorCode int) error {
	resBytes, err := json.Marshal(ErrorResponse{Success: false, Error: errorString})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	fmt.Fprintf(w, string(resBytes))
	if err != nil {
		return err
	}
	return nil
}

func GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *HandlerContext, hc *HomeAutoClaims) {
	u := UserResponse{
		Success: true,
		Message: "You're authenticated",
		Email:   hc.Email,
		Key:     hc.Prefix,
	}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		SendErrorResponse(w, "Problem parsing user info json", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonBytes))

}

func Auth(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *HandlerContext) {
	SetCorsHeaders(w)
	u, err := decodeUserFromRequest(r)
	if err != nil {
		SendErrorResponse(w, "Error: could not decode user. Did you POST with the proper user format? Full Error:"+err.Error(), 400)
		return
	}
	session, err := mgo.Dial(secrets.DB_DIAL_URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("homeauto").C("users")

	candidate := models.User{}
	c.Find(bson.M{"email": u.Email}).One(&candidate)
	berr := bcrypt.CompareHashAndPassword([]byte(candidate.Password), []byte(u.Password))
	if berr != nil {
		SendErrorResponse(w, berr.Error(), 400)
		return
	} else {
		claims := HomeAutoClaims{
			candidate.Email,
			candidate.Prefix,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * JWT_TTL).Unix(),
				Issuer:    "homeauto",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, jerr := token.SignedString(SecretKey)
		if jerr != nil {
			panic(jerr)
		}
		w.Header().Set("Content-Type", "application/json")
		resBytes, _ := json.Marshal(TokenResponse{Success: true, Token: tokenString})
		fmt.Fprintf(w, string(resBytes))
	}
}