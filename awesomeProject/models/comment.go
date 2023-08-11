package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type Comment struct {
	ID     int    `json:"id"`
	PostID int    `json:"postId"`
	UserID int    `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func InsertComment(db *gorm.DB, comment Comment) error {
	result := db.Create(&comment)
	if result.Error != nil {
		// If the comment already exists, we can ignore the error (duplicate key)
		if result.Error.Error() == "Error 1062: Duplicate entry" {
			return nil
		}
		return result.Error
	}
	return nil
}

func GetCommentsForPost(postID int) ([]Comment, error) {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?postId=%d", postID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var comments []Comment
	err = json.NewDecoder(resp.Body).Decode(&comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
