package api

import (
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"WASAPhoto/service/database"
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
	Owner         int       `json:"owner-id"` // UserID
	Username      string    `json:"owner-username"`
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
	CommentID     int       `json:"comment-id"`
	Timestamp     time.Time `json:"timestamp"`
	OwnerID       int       `json:"owner-id"`
	OwnerUsername string    `json:"owner-username"`
	Text          string    `json:"text"`
}

func (c *Comment) FromDatabase(comment database.Comment) {
	c.CommentID = comment.CommentID
	c.Timestamp = comment.Timestamp
	c.OwnerID = comment.UserID
	c.Text = comment.Text
}

// Array containing the images' filenames from the /images/ directory
var Images []string

const dirPath = "/tmp/images/"

var dirMutex = &sync.Mutex{}

func AddImage(data []byte, ext string) (string, error) {
	dirMutex.Lock()
	defer dirMutex.Unlock()

	newImageID := len(Images) + 1
	var filename = dirPath + "image" + strconv.Itoa(newImageID) + ext

	Images = append(Images, filename)

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	return "image" + strconv.Itoa(newImageID) + ext, nil
}
