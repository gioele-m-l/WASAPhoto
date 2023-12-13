package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-ype", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		ctx.Logger.Error("Content-Type header is not application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var username Username
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if it's a valid username
	err = username.checkUsername()
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid username")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := User{
		UserID:           -1,
		Username:         username.Username,
		ProfileImagePath: "",
		FollowersCount:   0,
		FollowingsCount:  0,
		PhotosCount:      0,
	}

	// Try to create the user in the db, if already exists return the token
	err = rt.db.CreateUser(user.ToDatabase())
	if err != nil {
		// The user is already in the database
		userDB, err := rt.db.GetUserByUsername(user.Username)
		if err != nil {
			ctx.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tokenDB, err := rt.db.GetUserToken(userDB.UserID)
		if err != nil {
			ctx.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var userToken UserToken
		userToken.FromDatabase(tokenDB)
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(userToken)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error in response")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	userDB, err := rt.db.GetUserByUsername(user.Username)
	if err != nil {
		ctx.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := getMD5Hash(user.Username + strconv.Itoa(userDB.UserID))
	userToken := UserToken{
		UserID: userDB.UserID,
		Token:  token,
	}
	err = rt.db.CreateToken(userToken.ToDatabase())
	if err != nil {
		ctx.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(userToken)
	if err != nil {
		ctx.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getMD5Hash(username string) string {
	hash := md5.Sum([]byte(username))
	return hex.EncodeToString(hash[:])
}
