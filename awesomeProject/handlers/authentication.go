package handlers

import (
	"awesomeProject/middlewares"
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

// @tags Authentication
// @Summary Register a new user
// @Description Register a new user with the provided data
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 200 {object} models.User
// @Router /register [post]
func RegisterHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		err := c.Bind(&user)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Check if the user already exists in the database
		var existingUser models.User
		res := db.Where("username = ?", user.Username).First(&existingUser)
		if res.RowsAffected > 0 {
			return c.String(http.StatusBadRequest, "User already exists")
		}

		// Hash the user's password and store it in the database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		user.Password = string(hashedPassword)

		// Insert the user into the database
		result := db.Create(&user)
		if result.Error != nil {
			return c.String(http.StatusInternalServerError, result.Error.Error())
		}

		return utils.JsonResponse(c, user)
	}
}

// @tags Authentication
// @Summary Log in as a user
// @Description Log in using the provided username and password
// @Accept json
// @Produce json
// @Param input body models.LoginInput true "Login data"
// @Success 200 {object} models.TokenResponse
// @Router /login [post]
func LoginHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// Bind the JSON input data to the 'input' struct
		err := c.Bind(&input)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Query the database to find the user with the provided username
		var user models.User
		result := db.Where("username = ?", input.Username).First(&user)
		if result.Error != nil {
			return c.String(http.StatusUnauthorized, "Incorrect username")
		}

		// Compare the provided password with the hashed password from the database
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			fmt.Println(user.Password, input.Password)
			return c.String(http.StatusUnauthorized, "Incorrect password")
		}

		// Create a JWT token for the authenticated user
		tokenString, err := middlewares.CreateJWTToken(user.ID, user.Username)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return utils.JsonResponse(c, map[string]string{
			"token": tokenString,
		})
	}
}
