package database

import "database/sql"

func (db *appdbimpl) UploadPhoto(caption string, pathToImage string, owner int) (int, error) {
	result, err := db.c.Exec(`INSERT INTO Photos (caption, path_to_image, owner) VALUES (?, ?, ?)`, caption, pathToImage, owner)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	return int(id), err
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
		if err == sql.ErrNoRows {
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
	return comments, nil
}

// Get likes by photo id
func (db *appdbimpl) GetLikesByPhotoID(photoID int) ([]int, error) {
	rows, err := db.c.Query(`SELECT likerID FROM Likes WHERE photoID = ?`, photoID)
	if err != nil {
		if err == sql.ErrNoRows {
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
	return likes, nil
}