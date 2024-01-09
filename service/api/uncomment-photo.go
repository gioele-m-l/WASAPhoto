package api

import (
	"WASAPhoto/service/api/reqcontext"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check the authentication
	userTok, err := CheckAuthentication(rt, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the photo-id from path
	photoID, err := strconv.Atoi(ps.ByName("photo-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: error in converting photo-id from string to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the comment-id froom path
	commentID, err := strconv.Atoi(ps.ByName("comment-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: error in converting comment-id from string to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user is the owner of the comment or if the user is the owner of the photo, then remove the comment
	err = rt.db.UncommentPhoto(photoID, commentID, userTok.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "Photo does not exists") {
			ctx.Logger.WithError(err).Info("uncommentPhoto")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "Comment does not exists under this photo") {
			ctx.Logger.WithError(err).Info("uncommentPhoto")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "User does not own the photo nor the comment") {
			ctx.Logger.WithError(err).Info("uncommentPhoto")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err).Error("uncommentPhoto: error in db query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the status code
	w.WriteHeader(http.StatusNoContent)

}
