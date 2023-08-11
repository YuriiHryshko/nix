package middlewares

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Retrieve the token from the request header
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.String(http.StatusUnauthorized, "Missing authorization token")
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			return c.String(http.StatusUnauthorized, "Invalid token")
		}

		// Extract user information from the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.String(http.StatusUnauthorized, "Invalid token data")
		}
		userID := uint(claims["id"].(float64))
		username := claims["username"].(string)

		// Attach user data to the request context
		c.Set("userID", userID)
		c.Set("username", username)

		return next(c)
	}
}

func CreateJWTToken(userID uint, username string) (string, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["username"] = username

	// Replace "your-secret-key" with your actual secret key for signing tokens
	secretKey := os.Getenv("JWT_SECRET_KEY")

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
