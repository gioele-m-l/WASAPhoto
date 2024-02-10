package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"WASAPhoto/service/api/reqcontext"
	"WASAPhoto/service/database"

	"github.com/julienschmidt/httprouter"
)

type response struct {
	Page   int     `json:"page"`
	Photos []Photo `json:"photos"`
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if there is an authorization token and if it's a valid token retrieving the related userID
	uTok, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in getMyStream function")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check the optional parameters 'page'
	// 'page' refers to the page number, in one page there are only a limited number of photos, incrementing the page number returns other photos(older), min 0 max 100
	var page int
	if r.URL.Query().Has("page") {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			ctx.Logger.WithError(err).Error("error in getMyStream function: cannot convert string to int")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		page = 0
	}
	if page < 0 {
		page = 0
	} else if page > 100 {
		page = 100
	}

	// Retrieve the user's stream, n photos in reverse chronological order, from followings that haven't blocked it
	var photosDB []database.Photo
	photosDB, err = rt.db.GetUserStream(uTok.UserID, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in getMyStream function: unable to retrieve user's stream from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Convert database.Photos in Photos objects
	var photos []Photo

	for _, p := range photosDB {
		var photo Photo
		photo.FromDatabase(p)

		user, err := rt.db.GetUserByID(photo.Owner)
		if err != nil {
			ctx.Logger.WithError(err).Error("error in getMyStream function: error retrieving user using photoID.owner")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		comments, err := rt.db.GetCommentsByPhotoID(p.PhotoID)
		if err != nil {
			ctx.Logger.WithError(err).Error("error in getMyStream function: error retrieving comments using photoID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		likes, err := rt.db.GetLikesByPhotoID(p.PhotoID)
		if err != nil {
			ctx.Logger.WithError(err).Error("error in getMyStream function: error retrieving likes using photoID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		photo.Username = user.Username
		photo.CommentsCount = len(comments)
		photo.LikesCount = len(likes)

		photos = append(photos, photo)
	}

	// Create a json object with the 'page' field and the photos array 'photos'
	w.Header().Set("Content-Type", "application/json")

	resp := response{
		Page:   page,
		Photos: photos,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in getMyStream function: error in response")
	}
}

func CheckAuthentication(rt *_router, r *http.Request) (database.UserToken, error) {
	var uTok database.UserToken
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return uTok, errors.New("missing authorization token")
	}
	uTok, err := rt.db.GetUserIDByAuthToken(authToken)
	if err != nil {
		return uTok, errors.New("invalid authorization token")
	}
	return uTok, nil
}
