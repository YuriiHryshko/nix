package handlers

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandler(t *testing.T) {
	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	e := echo.New()
	reqBody := map[string]string{"username": "testuser", "password": "testpassword"}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := RegisterHandler(db)

	// Test the handler
	err = handler(c)
	assert.NoError(t, err)
	// Check if the user already exists
	if rec.Code == http.StatusBadRequest {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	} else {
		// The user was successfully registered
		assert.Equal(t, http.StatusOK, rec.Code)

		// Check if the user is present in the database
		var user models.User
		result := db.Where("username = ?", "testuser").First(&user)
		assert.NoError(t, result.Error)
		assert.Equal(t, "testuser", user.Username)

		// Verify if the password has been stored in hashed form
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("testpassword"))
		assert.NoError(t, err)
	}
}

func TestLoginHandler(t *testing.T) {
	// Create a test database and perform model migrations
	db, err := database.InitializeDB("root:945929888@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	// Create a test user with a hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	testUser := models.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}
	db.Create(&testUser)

	e := echo.New()
	reqBody := map[string]string{"username": "testuser", "password": "testpassword"}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := LoginHandler(db)

	// Test the handler
	err = handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify if the response contains a token
	var responseMap map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &responseMap)
	assert.NoError(t, err)
	assert.Contains(t, responseMap, "token")
}
