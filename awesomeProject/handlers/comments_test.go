package handlers

import (
	"awesomeProject/database"
	"awesomeProject/middlewares"
	"awesomeProject/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommentHandler(t *testing.T) {
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
		"postId": 1,
		"name":   "John Doe",
		"email":  "johndoe@example.com",
		"body":   "This is a test comment.",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/comments", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenString) // Add the JWT token to the request header
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", testUser.ID) // Set the userID in the context (emulating the JWT middleware)
	handler := CreateCommentHandler(db)

	// Test the handler
	err = handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify if the response contains the created comment
	var responseComment models.Comment
	err = json.Unmarshal(rec.Body.Bytes(), &responseComment)
	assert.NoError(t, err)
	assert.Equal(t, 1, responseComment.PostID)
	assert.Equal(t, int(testUser.ID), responseComment.UserID)
	assert.Equal(t, "John Doe", responseComment.Name)
	assert.Equal(t, "johndoe@example.com", responseComment.Email)
	assert.Equal(t, "This is a test comment.", responseComment.Body)
}

func TestGetCommentsHandler_JSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/comments", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	handler := GetCommentsHandler(db)
	err = handler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json; charset=UTF-8", rec.Header().Get("Content-Type"))
	}
}

func TestGetCommentHandler_JSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/comments/1", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	// Prepare sample data if needed (create and insert sample comments)

	handler := GetCommentHandler(db)
	err = handler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json; charset=UTF-8", rec.Header().Get("Content-Type"))
	}
}

func TestUpdateCommentHandler(t *testing.T) {
	e := echo.New()
	reqJSON := `{"name": "Updated Name", "email": "updated@example.com", "body": "Updated Body"}`
	req := httptest.NewRequest(http.MethodPut, "/comments/1", strings.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	handler := UpdateCommentHandler(db)
	err = handler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json; charset=UTF-8", rec.Header().Get("Content-Type"))

		var updatedComment models.Comment
		db.First(&updatedComment, 1) // Assuming you are updating the comment with ID 1
		assert.Equal(t, "Updated Name", updatedComment.Name)
		assert.Equal(t, "updated@example.com", updatedComment.Email)
		assert.Equal(t, "Updated Body", updatedComment.Body)
	}
}

func TestDeleteCommentHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/comments/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	handler := DeleteCommentHandler(db)
	err = handler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Empty(t, rec.Body.String())

		// Verify if the comment was deleted from the database
		var deletedComment models.Comment
		result := db.Where("id = ?", 1).First(&deletedComment)
		assert.Error(t, result.Error)
	}
}
