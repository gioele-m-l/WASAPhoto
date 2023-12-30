package api

import (
	"net/http"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getImageFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check authorization header
	_, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Error("getImageFile function: missing or invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the file name from parameters
	imageID := ps.ByName("image-id")
	path := "/tmp/images/" + imageID
	file, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		ctx.Logger.WithError(err).Error("getImageFile function")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	file.Close()

	// Send the file
	w.Header().Set("Content-Type", "image/"+path[len(path)-3:])
	http.ServeFile(w, r, path)
}
