package api

import (
	"WASAPhoto/service/api/reqcontext"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check the authorization header
	userTok, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePhoto function: missing or invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the specified userID and check if it's related with the given token
	userID, err := strconv.Atoi(ps.ByName("user-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePhoto function: cannot convert string to int: user-id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if userID != userTok.UserID {
		ctx.Logger.Error("unlikePhoto function: forbidden")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Get the specified photoID and check if it's valid
	photoID, err := strconv.Atoi(ps.ByName("photo-id"))
	if err != nil {
		if err != nil {
			ctx.Logger.WithError(err).Error("unlikePhoto function: cannot convert string to int: photo-id")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	// Get the photo info
	photoDB, err := rt.db.GetPhotoByID(photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("unlikePhoto function: error in get photo query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if this user is banned by other user or viceversa
	banned, err := rt.db.CheckBan(photoDB.UserID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePhoto function: error in check ban query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	banned, err = rt.db.CheckBan(userID, photoDB.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePhoto function: error in check ban query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Remove the like from the photo
	_, err = rt.db.UnlikePhoto(photoID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePhoto function: db query error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Send back the response with the result of the operation
	w.WriteHeader(http.StatusNoContent)
}
