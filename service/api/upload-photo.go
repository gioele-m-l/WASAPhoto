package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"WASAPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if the user is logged in
	// Check if there is an authorization token
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		ctx.Logger.WithError(errors.New("missing user authorization token")).Error("authentication failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if it's valid
	userTok, err := rt.db.GetUserIDByAuthToken(authToken)
	if err != nil {
		ctx.Logger.WithError(err).Error("invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the content of the request
	file, header, err := r.FormFile("image")
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting the image uploaded")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filetype := header.Header.Get("Content-Type")
	var ext string
	switch filetype {
	case "image/png":
		ext = ".png"
	case "image/jpg":
		ext = ".jpg"
	default:
		ctx.Logger.Error("error in uploadProfileImage function: wrong mime type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	image, err := io.ReadAll(file)
	if err != nil {
		ctx.Logger.WithError(err).Error("error reading the file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	caption := r.FormValue("caption")

	// Check if the caption has max length 100
	if len(caption) > 100 {
		ctx.Logger.Error("error: caption is too long")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Storing image
	path, err := AddImage(image, ext)
	if err != nil {
		ctx.Logger.WithError(err).Error("error storing the image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Creating photo object
	photoID, err := rt.db.UploadPhoto(caption, path, userTok.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error inserting the photo in the db")
		err = os.Remove(path)
		if err != nil {
			ctx.Logger.WithError(err).Error("error removing the image")
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Returning the photo object
	var photo Photo
	photoDB, err := rt.db.GetPhotoByID(photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error retrieving photo using photoID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	photo.FromDatabase(photoDB)

	comments, err := rt.db.GetCommentsByPhotoID(photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error retrieving comments using photoID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	likes, err := rt.db.GetLikesByPhotoID(photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error retrieving likes using photoID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	photo.CommentsCount = len(comments)
	photo.LikesCount = len(likes)

	// Returning the json with the photo object
	w.Header().Set("Content-ype", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		ctx.Logger.WithError(err).Error("error sending the response")
	}
}
