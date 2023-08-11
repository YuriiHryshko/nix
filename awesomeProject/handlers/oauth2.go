package handlers

import (
	"awesomeProject/middlewares"
	"awesomeProject/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

func HandleGoogleAuth(c echo.Context) error {
	// Get the URL for Google OAuth2 authentication
	url := middlewares.GoogleOauthConfig.AuthCodeURL(middlewares.OauthStateString)
	return c.Redirect(http.StatusSeeOther, url)
}

func HandleGoogleCallback(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")

		// Exchange the authorization code for an access token
		token, err := middlewares.GoogleOauthConfig.Exchange(c.Request().Context(), code)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Use the access token to fetch user information
		client := middlewares.GoogleOauthConfig.Client(c.Request().Context(), token)
		response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer response.Body.Close()

		// Decode the response to get user profile
		var profile struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		}
		err = json.NewDecoder(response.Body).Decode(&profile)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Check if the user already exists in the database
		existingUser := models.User{}
		result := db.Where("email = ?", profile.Email).First(&existingUser)
		if result.RowsAffected == 0 {
			// User doesn't exist, register the user
			newUser := models.User{
				Email:    profile.Email,
				Password: "",
			}

			//  Insert the user into the database
			if err := db.Create(&newUser).Error; err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			// Create a JWT token for the new user
			newToken, er := middlewares.CreateJWTToken(newUser.ID, newUser.Username)
			if err != nil {
				return c.String(http.StatusInternalServerError, er.Error())
			}

			// Return the token to the user
			return c.JSON(http.StatusOK, map[string]string{
				"token":   newToken,
				"message": "New user registered and logged in via Google OAuth2.0",
			})
		}
		// Create a JWT token for the new user
		newToken, er := middlewares.CreateJWTToken(existingUser.ID, existingUser.Username)
		if err != nil {
			return c.String(http.StatusInternalServerError, er.Error())
		}

		// Return the token to the user
		return c.JSON(http.StatusOK, map[string]string{
			"token":   newToken,
			"message": "Existing user logged in via Google OAuth2.0",
		})
	}
}
