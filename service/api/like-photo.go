package api

import (
	"WASAPhoto/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check the authorization header
	userTok, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Error("likePhoto function: missing or invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the specified userID and check if it's related with the given token
	userID, err := strconv.Atoi(ps.ByName("user-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("likePhoto function: cannot convert string to int: user-id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if userID != userTok.UserID {
		ctx.Logger.Error("likePhoto function: forbidden")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Get the specified photoID and check if it's valid
	photoID, err := strconv.Atoi(ps.ByName("photo-id"))
	if err != nil {
		if err != nil {
			ctx.Logger.WithError(err).Error("likePhoto function: cannot convert string to int: photo-id")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// If the user isn't banned by the owner of the photo add a like to the photo
	rowsAffected, err := rt.db.LikePhoto(photoID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("likePhoto function: db query error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Send back the response with the result of the operation
	w.WriteHeader(http.StatusNoContent)
}
