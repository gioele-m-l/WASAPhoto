package api

import (
	"errors"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

// Username schema
type Username struct {
	Username string `json:"username"`
}

func (username *Username) checkUsername() error {
	if len(username.Username) < 3 || len(username.Username) > 16 {
		return errors.New("invalid username")
	}
	return nil
}

// User schema
type User struct {
	UserID           int    `json:"user-id"`
	Username         string `json:"username"`
	FollowersCount   int    `json:"followers-count"`
	FollowingsCount  int    `json:"followings-count"`
	PhotosCount      int    `json:"photos-count"`
	ProfileImagePath string `json:"profile-image-path"`
}

func (u *User) FromDatabase(user database.User) {
	u.UserID = user.UserID
	u.Username = user.Username
	u.ProfileImagePath = user.PathToProfileImage
}

func (u *User) ToDatabase() database.User {
	return database.User{
		UserID:             u.UserID,
		Username:           u.Username,
		PathToProfileImage: u.ProfileImagePath,
	}
}

// User-token struct
type UserToken struct {
	UserID int    `json:"user-id"`
	Token  string `json:"auth-token"`
}

func (ut *UserToken) FromDatabase(userToken database.UserToken) {
	ut.UserID = userToken.UserID
	ut.Token = userToken.Token
}

func (ut *UserToken) ToDatabase() database.UserToken {
	return database.UserToken{
		UserID: ut.UserID,
		Token:  ut.Token,
	}
}

// User summary
type UserSummary struct {
	UserID           int    `json:"user-id"`
	Username         string `json:"username"`
	ProfileImagePath string `json:"profile-image-path"`
}

// Image schema
type Image struct {
	File []byte
}

// Photo schema
type Photo struct {
	PhotoID       int       `json:"photo-id"`
	Timestamp     time.Time `json:"timestamp"`
	Owner         int       `json:"owner"` // UserID
	ImagePath     string    `json:"image-path"`
	LikesCount    int       `json:"likes-count"`
	CommentsCount int       `json:"comments-count"`
	Caption       string    `json:"caption"`
}

// Comment schema
type Comment struct {
	CommentID int       `json:"comment-id"`
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
}
