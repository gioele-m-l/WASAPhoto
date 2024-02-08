package api

import (
	"WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getPhotoLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check the authorization header
	_, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Error("getPhotoLike function: missing or invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the specified userID
	userID, err := strconv.Atoi(ps.ByName("user-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("getPhotoLike function: cannot convert string to int: user-id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the specified photoID and check if it's valid
	photoID, err := strconv.Atoi(ps.ByName("photo-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("getPhotoLike function: cannot convert string to int: photo-id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// If the user isn't banned by the owner of the photo add a like to the photo
	likers, err := rt.db.GetLikesByPhotoID(photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("getPhotoLike function: db query error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(likers); i++ {
		if likers[i] == userID {
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(userID)
			if err != nil {
				ctx.Logger.WithError(err).Error("getPhotoLike function: error in response")
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
