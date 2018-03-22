package handlers

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/suyashkumar/auth"
	"github.com/suyashkumar/conduit/server/db"
	"github.com/suyashkumar/conduit/server/device"
	"github.com/suyashkumar/conduit/server/entities"
)

func Hello(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	d device.Handler,
	db db.DatabaseHandler,
	a auth.Authenticator) {

	d.Call("suyash", "a", "ledToggle")
	fmt.Fprintf(w, "Hello, world")
}

func Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params, d device.Handler, db db.DatabaseHandler, a auth.Authenticator) {
	req := entities.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	// TODO: req validation
	if err != nil {
		logrus.WithError(err).Error("Could not parse RegisterRequest")
		err := sendJSON(w, entities.ErrorResponse{Error: "Could not parse RegisterRequest"}, 400)
		if err != nil {
			logrus.WithError(err).Error("!!!! Could not send error JSON response (RegisterRequest)")
		}
		return
	}

	// Create new user:
	u := auth.User{
		Email:              req.Email,
		MaxPermissionLevel: auth.PERMISSIONS_USER,
	}
	a.Register(&u, req.Password)

	// Create and add user's initial device secret
	logrus.Info(u.UUID)
	err = db.InsertDeviceSecret(u.UUID, entities.DeviceSecret{
		UUID:     uuid.NewV4(),
		UserUUID: u.UUID,
		Secret:   "hello",
	})

	if err != nil {
		logrus.WithError(err).WithField("user_uuid", u.UUID).Error("Error upserting device secret")
	}

	sendOK(w)
}

func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, d device.Handler, db db.DatabaseHandler, a auth.Authenticator) {
	req := entities.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	// TODO: req validation
	if err != nil {
		logrus.WithError(err).Error("Could not parse LoginRequest")
		err := sendJSON(w, entities.ErrorResponse{Error: "Could not parse LoginRequest"}, 400)
		if err != nil {
			logrus.WithError(err).Error("!!!! Could not send error JSON response (LoginRequest)")
		}
		return
	}

	// Get user if exists
	user, err := db.GetUser(auth.User{Email: req.Email})
	if err != nil {
		logrus.WithError(err).Error("Trouble fetching user")
		err := sendJSON(w, entities.ErrorResponse{Error: "Trouble fetching user"}, 400)
		if err != nil {
			logrus.WithError(err).Error("!!!! Could not send error JSON response (Login)")
		}
		return
	}

	// Get user's DeviceSecret to embed into Token
	secret, err := db.GetDeviceSecret(user.UUID)
	if err != nil {
		logrus.WithError(err).WithField("user_uuid", user.UUID).Error("Issue fetching device secret")
	}

	// Get Token for user
	token, err := a.GetToken(req.Email, req.Password, &auth.GetTokenOpts{
		RequestedPermissions: auth.PERMISSIONS_USER,
		Data:                 auth.TokenData{"deviceSecret": secret.Secret},
	})

	if err != nil {
		logrus.WithError(err).Error("Error getting token for user")
	}

	res := entities.LoginResponse{Token: token}
	sendJSON(w, res, 200)
}
