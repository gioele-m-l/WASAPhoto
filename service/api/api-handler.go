package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/login", rt.wrap(rt.doLogin))
	rt.router.PUT("/users/:username/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:username/profile-image", rt.wrap(rt.uploadProfileImage))
	rt.router.GET("/users/", rt.wrap(rt.listUsers))
	rt.router.GET("/users/:username/photos/", rt.wrap(rt.getUserPhotos))

	rt.router.GET("/users/:username/", rt.wrap(rt.getUserProfile))
	// rt.router.GET("/users/:username", rt.wrap(rt.getUserProfile))

	rt.router.PUT("/users/:username/followings/:user-id", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/followings/:user-id", rt.wrap(rt.unfollowUser))

	rt.router.PUT("/users/:username/banned/:user-id", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/banned/:user-id", rt.wrap(rt.unbanUser))

	rt.router.POST("/photos/", rt.wrap(rt.uploadPhoto))
	rt.router.GET("/photos/", rt.wrap(rt.getMyStream))
	rt.router.DELETE("/photos/:photo-id/", rt.wrap(rt.deletePhoto))
	rt.router.GET("/images/:image-id", rt.wrap(rt.getImageFile))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
