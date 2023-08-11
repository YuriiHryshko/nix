package middlewares

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

var (
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	// Random string for state parameter to prevent CSRF attacks
	OauthStateString = os.Getenv("OAUTH_STATE_STRING")
)
