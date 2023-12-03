package database

import (
	"database/sql"
	"errors"
)

// CreateUser insert a new user into the database
func (db *appdbimpl) CreateUser(u User) error {
	_, err := db.c.Exec(`INSERT INTO Users (username, path_to_profile_image) VALUES (?, "")`, u.Username)
	return err
}

// CreateToken insert a new token for an existing user into the database
func (db *appdbimpl) CreateToken(ut UserToken) error {
	_, err := db.c.Exec(`INSERT INTO AuthTokens (userID, token) VALUES (?, ?)`, ut.UserID, ut.Token)
	return err
}

// GetUser retrieve an existing user from the database with the username
func (db *appdbimpl) GetUserByUsername(username string) (User, error) {
	var user User
	err := db.c.QueryRow(`SELECT * FROM Users WHERE username = ?`, username).Scan(&user.UserID, &user.Username, &user.PathToProfileImage)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("User does not exists")
		}
	}
	return user, nil
}

// GetUserToken retrieve an existing user's token from the database using the user's usedID
func (db *appdbimpl) GetUserToken(userID int) (UserToken, error) {
	var userTok UserToken
	err := db.c.QueryRow(`SELECT * FROM AuthTokens WHERE userID = ?`, userID).Scan(&userTok.UserID, &userTok.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return userTok, errors.New("User's token does not exists")
		}
	}
	return userTok, nil
}

// SetMyUserName allow to change one user's username if it doesn't exists
func (db *appdbimpl) SetMyUserName(userID int, username string) error {
	_, err := db.c.Exec(`UPDATE Users SET username = ? WHERE userID = ?`, username, userID)
	return err
}
