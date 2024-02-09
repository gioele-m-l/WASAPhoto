package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"WASAPhoto/service/api/reqcontext"
	"WASAPhoto/service/database"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) listComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Check the authentication
	userTok, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Info("listComments: unauthorized access attemp")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the 'page' parameter if present
	/*
		var page int
		if r.URL.Query().Has("page") {
			page, err = strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				ctx.Logger.WithError(err).Error("listComments: error: cannot convert string to int")
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
	*/
	// Get and check the 'photo-id' parameter
	photoID, err := strconv.Atoi(ps.ByName("photo-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("listComments: error: cannot convert string to int (2)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get photo informations
	photoDB, err := rt.db.GetPhotoByID(photoID)
	if err != nil {
		ctx.Logger.WithError(err).Info("listComments: photo does not exists")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get photo owner and check the bans
	ownerID := photoDB.UserID
	banned, err := rt.db.CheckBan(ownerID, userTok.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("listComments: db query error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	banned, err = rt.db.CheckBan(userTok.UserID, ownerID)
	if err != nil {
		ctx.Logger.WithError(err).Error("listComments: db query error (2)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get the list of comments
	var commentsDB []database.Comment
	commentsDB, err = rt.db.GetCommentsByPhotoID(photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("listComments: db query error (3)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Conevrting database.Comment in Comment
	var comments []Comment
	for _, c := range commentsDB {
		var comm Comment
		comm.FromDatabase(c)
		commentOwner, err := rt.db.GetUserByID(comm.OwnerID)
		if err != nil {
			ctx.Logger.WithError(err).Error("listComments: db query error (4)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		comm.OwnerUsername = commentOwner.Username
		comments = append(comments, comm)
	}

	// Return the list of comments
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		ctx.Logger.WithError(err).Error("listComments: error sending the response back")
	}
}
