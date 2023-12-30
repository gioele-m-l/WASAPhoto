package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check the authorization header
	userTok, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in getUserPhotos: unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the username value and the userID
	username := ps.ByName("username")
	o_user, err := rt.db.GetUserByUsername(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in getUserPhotos: user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the user that made the request isn't blocked by the specified user
	blocked, err := rt.db.CheckBan(o_user.UserID, userTok.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in getUserPhotos")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if blocked {
		ctx.Logger.Info("getUserPhotos: the user is blocked by the other user")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Check if the user that made the request didn't block the other user
	blocked, err = rt.db.CheckBan(userTok.UserID, o_user.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in getUserPhotos")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if blocked {
		ctx.Logger.Info("getUserPhotos: the user blocked the other user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check the 'page' number
	var page int
	if r.URL.Query().Has("page") {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			ctx.Logger.WithError(err).Error("error in getUserPhotos function: cannot convert string to int")
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

	// Retrieve the photo objects of the user (20)
	var photosDB []database.Photo
	photosDB, err = rt.db.GetUserPhotos(o_user.UserID, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in getUserPhotos function: unable to retrieve user's stream from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert database.Photos in Photos objects
	var photos []Photo

	for _, p := range photosDB {
		var photo Photo
		photo.FromDatabase(p)

		comments, err := rt.db.GetCommentsByPhotoID(p.PhotoID)
		if err != nil {
			ctx.Logger.WithError(err).Error("error in getUserPhotos function: error retrieving comments using photoID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		likes, err := rt.db.GetLikesByPhotoID(p.PhotoID)
		if err != nil {
			ctx.Logger.WithError(err).Error("error in getUserPhotos function: error retrieving likes using photoID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		photo.CommentsCount = len(comments)
		photo.LikesCount = len(likes)

		photos = append(photos, photo)
	}

	// Send the response json back
	err = json.NewEncoder(w).Encode(photos)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in getUserPhotos function: error in response")
	}

}
