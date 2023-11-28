package database

import (
	"time"
)

type User struct {
	UserID             int
	Username           string
	PathToProfileImage string
}

type UserToken struct {
	UserID int
	Token  string
}

type Photo struct {
	PhotoID     int
	Timestamp   time.Time
	Caption     string
	PathToImage string
	UserID      int
}

type Comment struct {
	CommentID int
	Timestamp time.Time
	Text      string
	UserID    int // UserID del proprietario del commento
	PhotoID   int // PhotoID della foto commentata
}

type Image struct {
	Filename string
	Binary   []byte
}
