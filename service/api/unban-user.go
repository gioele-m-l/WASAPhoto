package api

import (
	"errors"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if the authentication token exists
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		ctx.Logger.WithError(errors.New("missing user authorization token")).Error("authentication failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userTok, err := rt.db.GetUserIDByAuthToken(authToken)
	if err != nil {
		ctx.Logger.WithError(err).Error("invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the username is correct
	var username Username
	username.Username = ps.ByName("username")
	err = username.checkUsername()
	if err != nil {
		ctx.Logger.WithError(err).Error("error in unbanUser function: invalid username")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user with that username exists
	userDB, err := rt.db.GetUserByUsername(username.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in unbanUser function: user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the user and the authentication token are related
	if userDB.UserID != userTok.UserID {
		ctx.Logger.Error("error in unbanUser function: forbidden")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check the userID value provided
	userID2, err := strconv.Atoi(ps.ByName("user-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("error in unbanUser function: error converting string to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the userID provided is the same of this user
	if userID2 == userDB.UserID {
		ctx.Logger.Warning("unbanUser function: user cannot unban itself")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if the user exists in db
	userDB2, err := rt.db.GetUserByID(userID2)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in unbanUser function: specified userID not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the relationship already exists
	e, err := rt.db.CheckBan(userDB.UserID, userDB2.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unbanUser: internal server errror")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !e {
		ctx.Logger.Info("unbanUser: the relationship already does not exist")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Delete the existent relationship
	err = rt.db.UnbanUser(userDB.UserID, userDB2.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unbanUser: internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
