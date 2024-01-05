package api

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strconv"

	"WASAPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check the authorization header
	userTok, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Error("deletePhoto function: missing or invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the photo id from path and check it
	r_photoID, err := strconv.Atoi(ps.ByName("photo-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("deletePhoto function")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the photo infos from db
	photoDB, err := rt.db.GetPhotoByID(r_photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("deletePhoto function: photo not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("deletePhoto function")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the owner of the specified photo is this user
	if photoDB.UserID != userTok.UserID {
		ctx.Logger.Error("deletePhoto function: user tried to delete another user's photo")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Delete the photo and the related image
	err = os.Remove(photoDB.PathToImage)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in deletePhoto function: cannot remove the image file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = rt.db.DeletePhoto(photoDB.PhotoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error in deletePhoto function: cannot delete photo from db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the status code (204)
	w.WriteHeader(http.StatusNoContent)
}
