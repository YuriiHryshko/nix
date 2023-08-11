package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func InsertPost(db *gorm.DB, post Post) error {
	result := db.Create(&post)
	if result.Error != nil {
		// If the post already exists, we can ignore the error (duplicate key)
		if result.Error.Error() == "Error 1062: Duplicate entry" {
			return nil
		}
		return result.Error
	}
	return nil
}

func GetPostsForUser(userID int) ([]Post, error) {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var posts []Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
