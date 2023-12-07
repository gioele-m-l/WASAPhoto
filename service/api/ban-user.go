package api

import (
	"errors"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if the authentication token exists
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		ctx.Logger.WithError(errors.New("Missing user authentication token")).Error("Authentication failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userTok, err := rt.db.GetUserIDByAuthToken(authToken)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the username is correct
	var username Username
	username.Username = ps.ByName("username")
	err = username.checkUsername()
	if err != nil {
		ctx.Logger.WithError(err).Error("Error in banUser function: invalid username")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user with that username exists
	userDB, err := rt.db.GetUserByUsername(username.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error in banUser function: user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the user and the authentication token are related
	if userDB.UserID != userTok.UserID {
		ctx.Logger.Error("Error in banUser function: forbidden")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check the userID value provided
	userID2, err := strconv.Atoi(ps.ByName("user-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("Error in banUser function: error converting string to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the userID provided is the same of this user
	if userID2 == userDB.UserID {
		ctx.Logger.Warning("banUser function: user cannot ban itself")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if the user exists in db
	userDB2, err := rt.db.GetUserByID(userID2)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error in banUser function: specified userID not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the other user has banned this user first
	e, err := rt.db.CheckBan(userDB2.UserID, userDB.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("banUser: internal server errror")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if e {
		ctx.Logger.Info("banUser: the other user has banned this user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the relationship already exists
	e, err = rt.db.CheckBan(userDB.UserID, userDB2.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("banUser: internal server errror")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if e {
		ctx.Logger.Info("banUser: the relationship already exists")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Create the new relationship
	err = rt.db.BanUser(userDB.UserID, userDB2.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("banUser: internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
