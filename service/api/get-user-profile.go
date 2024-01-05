package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"WASAPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Set the header for the response
	w.Header().Set("Content-Type", "application/json")

	// Check if the user is authenticated
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		ctx.Logger.WithError(errors.New("missing user authorization token")).Error("authentication failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err := rt.db.GetUserIDByAuthToken(authToken)
	if err != nil {
		ctx.Logger.WithError(err).Error("invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the username provided is valid and it's an existing user
	var username Username
	username.Username = ps.ByName("username")
	err = username.checkUsername()
	if err != nil {
		ctx.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Existing user ?
	userDB, err := rt.db.GetUserByUsername(username.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Creating the base User obj for the response
	var user User
	user.FromDatabase(userDB)

	// Retrievieng from the database other info of the user
	// Get followers
	user.FollowersCount, err = rt.db.GetUserFollowersCountByID(user.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in retrieving the followers count")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Get followings
	user.FollowingsCount, err = rt.db.GetUserFollowingsCountByID(user.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in retrieving the followings count")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Get the photos
	user.PhotosCount, err = rt.db.GetUserPhotosCountByID(user.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in retrieving the photos count")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Sending the response
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		ctx.Logger.WithError(err).Error("error sending the getUserProfile response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
