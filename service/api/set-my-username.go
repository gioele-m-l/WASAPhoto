package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"WASAPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "Application/json")
	// Check if the username format is correct
	var username Username
	username.Username = ps.ByName("username")
	err := username.checkUsername()
	if err != nil {
		ctx.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user with given username exists

	userDB, err := rt.db.GetUserByUsername(username.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var user User
	user.FromDatabase(userDB)

	// Check if the user is authenticated
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		ctx.Logger.WithError(errors.New("missing user authorization token")).Error("authentication failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userTokenDB, err := rt.db.GetUserToken(user.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("user token not found")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userTokenDB.Token != authToken {
		ctx.Logger.WithError(errors.New("forbidden")).Error("provided user token differ from the user's token in path")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		ctx.Logger.Error("Content-Type is not application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var username2 Username
	// Check the json provided
	err = json.NewDecoder(r.Body).Decode(&username2)
	if err != nil {
		ctx.Logger.WithError(err).Error("invalid json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check the username provided
	err = username2.checkUsername()
	if err != nil {
		ctx.Logger.WithError(err).Error("invalid username: ", username2.Username)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the username already exists
	_, err = rt.db.GetUserByUsername(username2.Username)
	if err != nil {
		err = rt.db.SetMyUserName(user.UserID, username2.Username)
		ctx.Logger.Info(user.UserID, username2.Username)
		if err != nil {
			ctx.Logger.WithError(err).Error("error in setting the username")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userDB, err = rt.db.GetUserByUsername(username2.Username)
		if err != nil {
			ctx.Logger.WithError(err).Error("error in getting back the modified user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userSum := UserSummary{
			UserID:           userDB.UserID,
			Username:         userDB.Username,
			ProfileImagePath: userDB.PathToProfileImage,
		}

		err = json.NewEncoder(w).Encode(userSum)
		if err != nil {
			ctx.Logger.WithError(err).Error("error in response")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	ctx.Logger.Info("user already exists")
	w.WriteHeader(http.StatusForbidden)
}
