package api

import (
	"errors"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if the authentication token exists
	authToken := r.Header.Get("Authentication")
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
		ctx.Logger.WithError(err).Error("Error in followUser function: invalid username")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user with that username exists
	userDB, err := rt.db.GetUserByUsername(username.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error in followUser function: user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the user and the authentication token are related
	if userDB.UserID != userTok.UserID {
		ctx.Logger.Error("Error in followUser function: forbidden")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check the userID value provided
	userID2, err := strconv.Atoi(ps.ByName("user-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("Error in followUser function: error converting string to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check if the userID provided is the same of this user
	if userID2 == userDB.UserID {
		ctx.Logger.Warning("followUser function: user cannot follow itself")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// Check if the user exists in db
	userDB2, err := rt.db.GetUserByID(userID2)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error in followUser function: specified userID not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if this user is blocked by the user with the specified userID
	// ...

	// Check if the relationship already exists
	err = rt.db.FollowUser(userDB.UserID, userDB2.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Warning("FollowUser: this relationship already exists")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
