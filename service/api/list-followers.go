package api

import (
	"WASAPhoto/service/api/reqcontext"
	"WASAPhoto/service/database"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) listFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if there's an authorization token and if it's valid
	_, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Error("listFollowers: missing or invalid authorization token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the 'username' parameter
	username := ps.ByName("username")
	if len(username) < 3 || len(username) > 16 {
		ctx.Logger.Error("listFollowers: invalid username")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err = rt.db.GetUserByUsername(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("listBanned: error retrieving the user with username")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var dbUsers []database.User
	dbUsers, err = rt.db.ListFollowers(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("listFollowers: cannot retrieve users from db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var users []UserSummary
	for _, userDB := range dbUsers {
		users = append(users, UserSummary{
			UserID:           userDB.UserID,
			Username:         userDB.Username,
			ProfileImagePath: userDB.PathToProfileImage,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}