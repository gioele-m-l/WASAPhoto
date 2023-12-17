/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// GetName() (string, error)
	// SetName(name string) error
	CreateUser(u User) error
	CreateToken(ut UserToken) error

	GetUserByUsername(username string) (User, error)
	GetUserByID(userID int) (User, error)

	GetUserToken(userID int) (UserToken, error)
	SetMyUserName(userID int, username string) error
	GetUserIDByAuthToken(token string) (UserToken, error)
	ListUsers(substring string) ([]User, error)

	GetUserFollowersCountByID(userID int) (int, error)
	GetUserFollowingsCountByID(userID int) (int, error)
	GetUserPhotosCountByID(userID int) (int, error)

	FollowUser(followerID int, followedID int) error
	UnfollowUser(followerID int, followedID int) error

	CheckBan(blockerID int, blockedID int) (bool, error)
	BanUser(blockerID int, blockedID int) error
	UnbanUser(blockerID int, blockedID int) error

	UpdateProfileImage(uID int, path string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Turning on the support for foreign keys
	_, errPragma := db.Exec(`PRAGMA foreign_keys= ON`)
	if errPragma != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPragma)
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		/*
			sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
			_, err = db.Exec(sqlStmt)
			if err != nil {
				return nil, fmt.Errorf("error creating database structure: %w", err)
			}
		*/

		// Creating the Users table
		usersTable := `CREATE TABLE IF NOT EXISTS Users (
			userID INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(16) UNIQUE NOT NULL,
			path_to_profile_image TEXT NOT NULL
		);`
		_, errUsersTable := db.Exec(usersTable)
		if errUsersTable != nil {
			return nil, fmt.Errorf("error creating the Users table in database structure: %w", errUsersTable)
		}

		// Creating the Followers table
		followersTable := `CREATE TABLE IF NOT EXISTS Followers (
			followedID INTEGER,
			followerID INTEGER,
			PRIMARY KEY(followedID, followerID),
			FOREIGN KEY(followedID) REFERENCES Users(userID) ON DELETE CASCADE,
			FOREIGN KEY(followerID) REFERENCES Users(userID) ON DELETE CASCADE
		);`
		_, errFollowersTable := db.Exec(followersTable)
		if errFollowersTable != nil {
			return nil, fmt.Errorf("error creating the Followers table in database structure: %w", errFollowersTable)
		}

		// Creating the Blocked_users table
		blockedUsersTable := `CREATE TABLE IF NOT EXISTS Blocked_users (
			blockerID INTEGER,
			blockedID INTEGER,
			PRIMARY KEY(blockerID, blockedID),
			FOREIGN KEY(blockerID) REFERENCES Users(userID) ON DELETE CASCADE,
			FOREIGN KEY(blockedID) REFERENCES Users(userID) ON DELETE CASCADE
		);`
		_, errBlockedUsersTable := db.Exec(blockedUsersTable)
		if errBlockedUsersTable != nil {
			return nil, fmt.Errorf("error creating the Blocked_users table in database structure: %w", errBlockedUsersTable)
		}

		// Creating the Photos table
		photosTable := `CREATE TABLE IF NOT EXISTS Photos (
			photoID INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			caption VARCHAR(100),
			path_to_image TEXT NOT NULL, 
			owner INTEGER,
			FOREIGN KEY(owner) REFERENCES Users(userID) ON DELETE CASCADE
		);`
		_, errPhotosTable := db.Exec(photosTable)
		if errPhotosTable != nil {
			return nil, fmt.Errorf("error creating the Photos table in database structure: %w", errUsersTable)
		}

		// Creating the Likes table
		likesTable := `CREATE TABLE IF NOT EXISTS Likes (
			photoID INTEGER,
			likerID INTEGER,
			PRIMARY KEY(photoID, likerID),
			FOREIGN KEY(photoID) REFERENCES Photos(photoID) ON DELETE CASCADE,
			FOREIGN KEY(likerID) REFERENCES Users(userID) ON DELETE CASCADE
		);`
		_, errLikesTable := db.Exec(likesTable)
		if errLikesTable != nil {
			return nil, fmt.Errorf("error creating the Likes table in database structure: %w", errLikesTable)
		}

		// Creating the Comments table
		commentsTable := `CREATE TABLE IF NOT EXISTS Comments (
			commentID INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			text VARCHAR(256),
			commenterID INTEGER,
			photoID INTEGER,
			FOREIGN KEY(commenterID) REFERENCES Users(userID) ON DELETE CASCADE,
			FOREIGN KEY(photoID) REFERENCES Photos(photoID) ON DELETE CASCADE
		);`
		_, errCommentsTable := db.Exec(commentsTable)
		if errCommentsTable != nil {
			return nil, fmt.Errorf("error creating the Comments table in database structure: %w", errCommentsTable)
		}

		// Creating the AuthTokens table
		authTokensTable := `CREATE TABLE IF NOT EXISTS AuthTokens (
			userID INTEGER PRIMARY KEY,
			token VARCHAR(32) NOT NULL,
			FOREIGN KEY(userID) REFERENCES Users(userID) ON DELETE CASCADE
		);`
		_, errAuthTokensTable := db.Exec(authTokensTable)
		if errAuthTokensTable != nil {
			return nil, fmt.Errorf("error creating the AuthTokens table in database structure: %w", errAuthTokensTable)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
