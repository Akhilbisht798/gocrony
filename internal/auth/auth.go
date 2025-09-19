package auth

import (
	"log"
	"net/http"

	"github.com/akhilbisht798/gocrony/config"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	MaxAge = 86400 * 30
	isProd = false
)

type Provider string
const (
	Email Provider = "email"
	Google Provider = "google"
)

func NewAuth() {
	googleClientId := config.GetEnv("GOOGLE_CLIENT_ID", "")
	googleClientSecret := config.GetEnv("GOOGLE_CLIENT_SECRET", "")
	googleCallbackURL := config.GetEnv("GOOGLE_URI_CALLBACK", "")

	sessionKey := config.GetEnv("SESSION_KEY", "secret")

	log.Println("Initializing Google OAuth provider")
	log.Println("Google Client ID:", googleClientId)
	log.Println("Callback URL:", googleCallbackURL)
	log.Println("Session key:", sessionKey)

	store := sessions.NewCookieStore([]byte(sessionKey))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	store.Options.SameSite = http.SameSiteLaxMode // Important for OAuth

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, googleCallbackURL),
	)
}
