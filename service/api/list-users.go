package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) listUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	substring := ""
	if r.URL.Query().Has("search") {
		substring = r.URL.Query().Get("search")
	}

	// Get the authentication token and check if it's a registered user
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

	// Get a list of user summaries from the db (limited to 100 results)
	var users []UserSummary
	usersDB, err := rt.db.ListUsers(substring)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in returning the users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, userDB := range usersDB {
		users = append(users, UserSummary{
			UserID:           userDB.UserID,
			Username:         userDB.Username,
			ProfileImagePath: userDB.PathToProfileImage,
		})
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
