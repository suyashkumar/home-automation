package handlers

import (
	"net/http"

	"encoding/json"

	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/suyashkumar/auth"
	"github.com/suyashkumar/conduit/db"
	"github.com/suyashkumar/conduit/device"
	"github.com/suyashkumar/conduit/entities"
	sec "github.com/suyashkumar/conduit/secret"
)

// Register allows a new user to create an account with Conduit
func Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params, d device.Handler, db db.Handler, a auth.Authenticator) {
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
	err = db.InsertAccountSecret(u.UUID, entities.AccountSecret{
		UUID:     uuid.NewV4(),
		UserUUID: u.UUID,
		Secret:   sec.GetRandString(10),
	})

	if err != nil {
		logrus.WithError(err).WithField("user_uuid", u.UUID).Error("Error upserting device secret")
	}

	sendOK(w)
}

// Login allows the user to authenticate with conduit and get a freshly minted JWT
func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, d device.Handler, db db.Handler, a auth.Authenticator) {
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

	// Get user's AccountSecret to embed into Token
	secret, err := db.GetAccountSecret(user.UUID)
	if err != nil {
		logrus.WithError(err).WithField("user_uuid", user.UUID).Error("Issue fetching device secret")
	}

	// Get Token for user
	token, err := a.GetToken(req.Email, req.Password, &auth.GetTokenOpts{
		RequestedPermissions: auth.PERMISSIONS_USER,
		Data:                 auth.TokenData{ACCOUNT_SECRET_KEY: secret.Secret},
	})

	if err != nil {
		logrus.WithError(err).Error("Error getting token for user")
	}

	res := entities.LoginResponse{Token: token}
	sendJSON(w, res, 200)
}

// Call allows a user to issue an RPC to one of their devices and optionally get a response from the device
func Call(w http.ResponseWriter, r *http.Request, ps httprouter.Params, d device.Handler, db db.Handler, a auth.Authenticator) {
	req := entities.CallRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.WithError(err).Error("Could not parse CallRequest")
		err := sendJSON(w, entities.ErrorResponse{Error: "Could not parse CallRequest"}, 400)
		if err != nil {
			logrus.WithError(err).Error("!!!! Could not send error JSON response (CallRequest)")
		}
		return
	}

	// Authenticate User
	claims, err := a.Validate(req.Token)
	if err == auth.ErrorValidatingToken {
		logrus.WithField("token", req.Token).Info("Error validating token")
		err := sendJSON(w, entities.ErrorResponse{Error: "Error validating token"}, 401)
		if err != nil {
			logrus.WithError(err).Error("!!!! Could not send error JSON response (CallRequest)")
		}
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Unknown error validating token")
		err := sendJSON(w, entities.ErrorResponse{Error: "Error validating token"}, 500)
		if err != nil {
			logrus.WithError(err).Error("!!!! Could not send error JSON response (CallRequest)")
		}
		return
	}

	c := d.Call(req.DeviceName, claims.Data[ACCOUNT_SECRET_KEY], req.FunctionName, req.WaitForDeviceResponse)

	if req.WaitForDeviceResponse {
		select {
		case res := <-c:
			logrus.WithField("response", res).Info("Device responded")
			r := entities.SendResponse{
				Response: res,
			}
			sendJSON(w, r, 200)
		case <-time.After(3 * time.Second):
			logrus.Warn("Timed out waiting for device response")
			e := entities.ErrorResponse{
				Error: "Timed out while waiting for the device to respond",
			}
			sendJSON(w, e, 500)
		}
	} else {
		sendOK(w)
	}
}

// UserInfo returns information to the user about their account, including their current account secret
func UserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, d device.Handler, db db.Handler, a auth.Authenticator) {
	req := entities.UserInfoRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.WithError(err).Error("Could not parse CallRequest")
		err := sendJSON(w, entities.ErrorResponse{Error: "Could not parse CallRequest"}, 400)
		if err != nil {
			logrus.WithError(err).Error("!!!! Could not send error JSON response (CallRequest)")
		}
		return
	}

	// Authenticate User. TODO: factor out
	claims, err := a.Validate(req.Token)
	if err == auth.ErrorValidatingToken {
		logrus.WithField("token", req.Token).Info("Error validating token")
		err := sendJSON(w, entities.ErrorResponse{Error: "Error validating token"}, 401)
		if err != nil {
			logrus.WithError(err).Error("!!!! Could not send error JSON response (UserInfoRequest)")
		}
		return
	}
	if err != nil {
		logrus.WithError(err).Error("Unknown error validating token")
		err := sendJSON(w, entities.ErrorResponse{Error: "Error validating token"}, 500)
		if err != nil {
			logrus.WithError(err).Error("!!!! Could not send error JSON response")
		}
		return
	}

	sendJSON(w, entities.UserInfoResponse{AccountSecret: claims.Data[ACCOUNT_SECRET_KEY]}, 200)

}
