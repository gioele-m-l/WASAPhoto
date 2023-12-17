package api

import (
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
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

func (p *Photo) FromDatabase(photo database.Photo) {
	p.PhotoID = photo.PhotoID
	p.Timestamp = photo.Timestamp
	p.Owner = photo.UserID
	p.ImagePath = photo.PathToImage
	p.Caption = photo.Caption
}

// Comment schema
type Comment struct {
	CommentID int       `json:"comment-id"`
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
}

// Array containing the images' filenames from the /images/ directory
var Images []string
var dirPath = "./images/"
var dirMutex = &sync.Mutex{}

func AddImage(data []byte, ext string, ctx reqcontext.RequestContext) (string, error) {
	dirMutex.Lock()
	defer dirMutex.Unlock()

	newImageID := len(Images) + 1
	var filename = dirPath + "image" + strconv.Itoa(newImageID) + ext

	Images = append(Images, filename)

	file, err := os.Create(filename)
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating the image file")
		return "", err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		ctx.Logger.WithError(err).Error("error writing the image file")
		return "", err
	}

	return filename, nil
}
