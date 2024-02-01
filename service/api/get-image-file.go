package api

import (
	"encoding/base64"
	"net/http"
	"os"

	"WASAPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getImageFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check authorization header
	ctx.Logger.Info("getImageFile: request received")
	_, err := CheckAuthentication(rt, r)
	if err != nil {
		ctx.Logger.WithError(err).Error("getImageFile function: missing or invalid user token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the file name from parameters
	imageID := ps.ByName("image-id")
	ctx.Logger.Info(imageID)
	path := "/tmp/images/" + imageID
	data, err := os.ReadFile(path)
	ctx.Logger.Info("Reading: " + path)
	if err != nil {
		ctx.Logger.WithError(err).Error("getImageFile function")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send the file
	w.Header().Set("Content-Type", "image/"+path[len(path)-3:])
	data = []byte(base64.StdEncoding.EncodeToString(data))
	ctx.Logger.Info("Serving: " + path)
	w.Write(data)
}
