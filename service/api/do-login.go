package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var username Username
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} // Add a check to the username

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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tokenDB, err := rt.db.GetUserToken(userDB.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var userToken UserToken
		userToken.FromDatabase(tokenDB)
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(userToken)
		return
	}

	userDB, err := rt.db.GetUserByUsername(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := getMD5Hash(user.Username)
	fmt.Println(token)
	userToken := UserToken{
		UserID: userDB.UserID,
		Token:  token,
	}
	err = rt.db.CreateToken(userToken.ToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userToken)
}

func getMD5Hash(username string) string {
	hash := md5.Sum([]byte(username))
	return hex.EncodeToString(hash[:])
}
