package handlers

import (
	"awesomeProject/database"
	"awesomeProject/middlewares"
	"awesomeProject/models"
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreatePostHandler(t *testing.T) {
	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	// Create a test user
	testUser := models.User{
		Username: "testuser",
		Password: "testpassword",
	}
	db.Create(&testUser)

	// Create a JWT token for the test user
	tokenString, err := middlewares.CreateJWTToken(testUser.ID, testUser.Username)
	assert.NoError(t, err)

	e := echo.New()
	reqBody := map[string]interface{}{
		"title": "New Post",
		"body":  "This is a new test post.",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", testUser.ID)
	handler := CreatePostHandler(db)

	err = handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responsePost models.Post
	err = json.Unmarshal(rec.Body.Bytes(), &responsePost)
	assert.NoError(t, err)
	assert.Equal(t, "New Post", responsePost.Title)
	assert.Equal(t, "This is a new test post.", responsePost.Body)
	assert.Equal(t, int(testUser.ID), responsePost.UserID)
}

func TestGetPostsHandler_JSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	handler := GetPostsHandler(db)
	err = handler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json; charset=UTF-8", rec.Header().Get("Content-Type"))
	}
}

func TestGetPostHandler_JSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	// Prepare sample data if needed (create and insert sample comments)

	handler := GetPostHandler(db)
	err = handler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json; charset=UTF-8", rec.Header().Get("Content-Type"))
	}
}

func TestUpdatePostHandler(t *testing.T) {
	e := echo.New()
	reqJSON := `{"title": "Updated Title", "body": "Updated Body"}`
	req := httptest.NewRequest(http.MethodPut, "/posts/1", strings.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	handler := UpdatePostHandler(db)
	err = handler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json; charset=UTF-8", rec.Header().Get("Content-Type"))

		var updatedPost models.Post
		db.First(&updatedPost, 1)
		assert.Equal(t, "Updated Title", updatedPost.Title)
		assert.Equal(t, "Updated Body", updatedPost.Body)
	}
}

func TestDeletePostHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	handler := DeletePostHandler(db)
	err = handler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Empty(t, rec.Body.String())

		// Verify if the comment was deleted from the database
		var deletedPost models.Post
		result := db.Where("id = ?", 1).First(&deletedPost)
		assert.Error(t, result.Error)
	}
}
