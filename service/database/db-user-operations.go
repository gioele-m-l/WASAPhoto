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

// Get the userID of the user by using its authentication token
func (db *appdbimpl) GetUserIDByAuthToken(token string) (UserToken, error) {
	var userTok UserToken
	err := db.c.QueryRow(`SELECT * FROM AuthTokens WHERE token = ?`, token).Scan(&userTok.UserID, &userTok.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return userTok, errors.New("User's token doesn't exists")
		}
	}
	return userTok, err
}

// Get a list of users (max 100)
func (db *appdbimpl) ListUsers(substring string) ([]User, error) {
	var users []User
	stmt, err := db.c.Prepare("SELECT * FROM Users WHERE username LIKE ? LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(substring + "%")
	defer rows.Close()

	for rows.Next() {
		var u User
		err := rows.Scan(&u.UserID, &u.Username, &u.PathToProfileImage)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	err = rows.Err()
	if err != nil {
		return users, err
	}

	return users, err
}

// Get the number of followers of the specified user
func (db *appdbimpl) GetUserFollowersCountByID(userID int) (int, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(followerID) FROM Followers WHERE followedID = ?`, userID).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// Get the number of followings of the specified user
func (db *appdbimpl) GetUserFollowingsCountByID(userID int) (int, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(followedID) FROM Followers WHERE followerID = ?`, userID).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// Get the number of photos of the specified user
func (db *appdbimpl) GetUserPhotosCountByID(userID int) (int, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(photoID) FROM Photos WHERE owner = ?`, userID).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// Get the user object specifying the userID
func (db *appdbimpl) GetUserByID(userID int) (User, error) {
	var user User
	err := db.c.QueryRow(`SELECT * FROM Users WHERE userID = ?`, userID).Scan(&user.UserID, &user.Username, &user.PathToProfileImage)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Insert a new relationship Follower in db
func (db *appdbimpl) FollowUser(followerID int, followedID int) error {
	_, err := db.c.Exec(`INSERT INTO Followers (followedID, followerID) VALUES (?, ?)`, followedID, followerID)
	if err != nil {
		return err
	}
	return nil
}
