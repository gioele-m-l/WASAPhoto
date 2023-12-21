package api

import (
	"errors"
	"io"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadProfileImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if the user provided the authorization token and if it's valid
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

	// Check if the the username in the path is valid and the relative user exists in db
	var username Username
	username.Username = ps.ByName("username")
	err = username.checkUsername()
	if err != nil {
		ctx.Logger.WithError(err).Error("error in uploadProfileImage function: invalid username")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userDB, err := rt.db.GetUserByUsername(username.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in uploadProfileImage function: user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the user (by username) and the token are related
	if userDB.UserID != userTok.UserID {
		ctx.Logger.Error("error in uploadProfileImage function: forbidden")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check the content-type of the request body, must be image/png or image/jpg
	contType := r.Header.Get("Content-Type")
	var ext string
	switch contType {
	case "image/png":
		ext = ".png"
	case "image/jpg":
		ext = ".jpg"
	default:
		ctx.Logger.Error("error in uploadProfileImage function: bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the file uploaded and save it in the images/ directory
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in uploadProfileImage function")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path, err := AddImage(data, ext)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in uploadProfileImage function: AddImage")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Update the user record to link it to its profile image

	err = rt.db.UpdateProfileImage(userDB.UserID, path)
	if err != nil {
		ctx.Logger.WithError(err).Error("error updating db informations")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
