package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) UploadPhoto(caption string, pathToImage string, owner int) (int, error) {
	result, err := db.c.Exec(`INSERT INTO Photos (caption, path_to_image, owner) VALUES (?, ?, ?)`, caption, pathToImage, owner)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (db *appdbimpl) DeletePhoto(photoID int) error {
	_, err := db.c.Exec(`DELETE FROM Photos WHERE photoID = ?`, photoID)
	return err
}

func (db *appdbimpl) GetPhotoByID(photoID int) (Photo, error) {
	var photo Photo
	err := db.c.QueryRow(`SELECT * FROM Photos WHERE photoID = ?`, photoID).Scan(&photo.PhotoID, &photo.Timestamp, &photo.Caption, &photo.PathToImage, &photo.UserID)
	return photo, err
}

// Get comments by photo ID
func (db *appdbimpl) GetCommentsByPhotoID(photoID int) ([]Comment, error) {
	rows, err := db.c.Query(`SELECT * FROM Comments WHERE photoID = ?`, photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.CommentID, &comment.PhotoID, &comment.UserID, &comment.Timestamp, &comment.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	err = rows.Err() // Check if there was an error during the iteration
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// Get likes by photo id
func (db *appdbimpl) GetLikesByPhotoID(photoID int) ([]int, error) {
	rows, err := db.c.Query(`SELECT likerID FROM Likes WHERE photoID = ?`, photoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var likes []int
	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		likes = append(likes, userID)
	}

	err = rows.Err() // Check if there was an error during the iteration
	if err != nil {
		return nil, err
	}

	return likes, nil
}

// Get user stream using userID
func (db *appdbimpl) GetUserStream(userID int, page int) ([]Photo, error) {
	var photos []Photo
	query := `SELECT Photos.* FROM Photos JOIN Followers ON Photos.owner = Followers.followedID
				WHERE Followers.followerID = ? ORDER BY Photos.timestamp DESC LIMIT ? OFFSET ?`

	limit := 20
	offset := 20 * page

	rows, err := db.c.Query(query, userID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return photos, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var photo Photo
		err := rows.Scan(&photo.PhotoID, &photo.Timestamp, &photo.Caption, &photo.PathToImage, &photo.UserID)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	err = rows.Err() // Check if there was an error during the iteration
	if err != nil {
		return nil, err
	}

	return photos, nil
}

// Get user photos using userID
func (db *appdbimpl) GetUserPhotos(userID int, page int) ([]Photo, error) {
	var photos []Photo
	query := `SELECT * FROM Photos WHERE Photos.owner = ? ORDER BY Photos.timestamp DESC LIMIT ? OFFSET ?`

	limit := 20
	offset := 20 * page

	rows, err := db.c.Query(query, userID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return photos, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var photo Photo
		err := rows.Scan(&photo.PhotoID, &photo.Timestamp, &photo.Caption, &photo.PathToImage, &photo.UserID)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	err = rows.Err() // Check if there was an error during the iteration
	if err != nil {
		return nil, err
	}

	return photos, nil
}

// Like photo
func (db *appdbimpl) LikePhoto(photoID int, userID int) (int64, error) {
	result, err := db.c.Exec(`INSERT INTO Likes (photoID, likerID) SELECT ?, ? WHERE EXISTS(
									SELECT 1 FROM Photos WHERE photoID = ?
								) AND NOT EXISTS(
									SELECT 1 FROM Blocked_users INNER JOIN Photos ON Blocked_users.blockerID = Photos.owner
									WHERE Blocked_users.blockedID = ? AND Photos.photoID = ?
								) AND NOT EXISTS(
									SELECT 1 FROM Blocked_users INNER JOIN Photos ON Blocked_users.blockedID = Photos.owner
									WHERE Blocked_users.blockerID = ? AND Photos.photoID = ?
								)`, photoID, userID, photoID, userID, photoID, userID, photoID, photoID, userID)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

// Unlike photo
func (db *appdbimpl) UnlikePhoto(photoID int, userID int) (int64, error) {
	result, err := db.c.Exec(`DELETE FROM Likes WHERE photoID = ? AND likerID = ?`, photoID, userID)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// Comment photo
func (db *appdbimpl) CommentPhoto(photoID int, userID int, commentText string) (sql.Result, error) {
	result, err := db.c.Exec(`INSERT INTO Comments (text, commenterID, photoID) SELECT ?, ?, ? WHERE EXISTS (
									SELECT 1 FROM Photos WHERE photoID = ?
								) AND NOT EXISTS (
									SELECT 1 FROM Blocked_users INNER JOIN Photos ON Blocked_users.blockerID = Photos.owner
									WHERE Blocked_users.blockedID = ? AND Photos.photoID = ?
								) AND NOT EXISTS(
									SELECT 1 FROM Blocked_users INNER JOIN Photos ON Blocked_users.blockedID = Photos.owner
									WHERE Blocked_users.blockerID = ? AND Photos.photoID = ?
								)`, commentText, userID, photoID, photoID, userID, photoID, userID, photoID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get comment by ID
func (db *appdbimpl) GetCommentByID(commentID int64) (Comment, error) {
	var comment Comment
	err := db.c.QueryRow(`SELECT * FROM Comments WHERE commentID = ?`, commentID).Scan(&comment.CommentID, &comment.Timestamp, &comment.Text, &comment.UserID, &comment.PhotoID)
	if err != nil {
		return comment, err
	}
	return comment, nil
}

// Uncomment photo
func (db *appdbimpl) UncommentPhoto(photoID int, commentID int, userID int) error {
	// Begin the transaction
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}

	// Check if the photo and the comment exist, then check if the comment belongs to the photo
	row := tx.QueryRow(`SELECT CASE 
						WHEN NOT EXISTS(
							SELECT 1 FROM Photos WHERE photoID = ?
						) THEN 1 
						WHEN NOT tEXISTS(
							SELECT 1 FROM Comments WHERE commentID = ? AND photoID = ?
						) THEN 2
						ELSE 0 END`, photoID, commentID, photoID)
	var code int
	err = row.Scan(&code)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return err2
		}
		return err
	}

	switch code {
	case 1:
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return fmt.Errorf("Photo does not exists")
	case 2:
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return fmt.Errorf("Comment does not exists under this photo")
	case 0:
		break
	}

	// Check if the photo or the comment belong to the user
	var count int
	err = tx.QueryRow(`SELECT COUNT(*) FROM Photos p LEFT JOIN Comments c ON p.photoID = c.photoID
						WHERE p.photoID = ? AND c.commentID = ? AND (p.owner = ? OR c.commenterID = ?)
						`, photoID, commentID, userID, userID).Scan(&count)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return err2
		}
		return err
	}

	if count == 0 {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return fmt.Errorf("User does not own the photo nor the comment")
	}

	// Delete the comment
	_, err = tx.Exec(`DELETE FROM Comments WHERE commentID = ?`, commentID)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return err2
		}
		return err
	}

	// Exec the transaction
	err = tx.Commit()
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return err2
		}
		return err
	}
	return nil
}
