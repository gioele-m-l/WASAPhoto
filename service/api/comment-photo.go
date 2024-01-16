package api

import (
	"WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check the authentication
	userTok, err := CheckAuthentication(rt, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the photo-id from path
	photoID, err := strconv.Atoi(ps.ByName("photo-id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error in converting photo-id from string to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the content of the body
	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error in parsing the request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check the content of the body
	if len(comment.Text) < 1 || len(comment.Text) > 255 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//  Check if this user is banned and viceversa and create the relationship i db
	result, err := rt.db.CommentPhoto(photoID, userTok.UserID, comment.Text)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error in db query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error in db query (2)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error in db query (3)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the comment
	commentDB, err := rt.db.GetCommentByID(commentID)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error in db query (4)")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	comment = Comment{
		CommentID: commentDB.CommentID,
		Timestamp: commentDB.Timestamp,
		Text:      commentDB.Text,
	}

	// Send back the response with the object created
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error sending back the response")
	}
}
