package config

import (
	"os"
	"time"
)

type Config struct {
	JWTSecret        string
	AccessTTL        time.Duration
	RefreshTTL       time.Duration
	RefreshCookieKey string
	CookieSecure     bool
	CorsOrigin       string
	Port             string
}

func MustLoad() Config {
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "dev-secret-change-me"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "http://localhost:5173"
	}
	return Config{
		JWTSecret:        sec,
		AccessTTL:        15 * time.Minute,
		RefreshTTL:       14 * 24 * time.Hour,
		RefreshCookieKey: "refreshToken",
		CookieSecure:     os.Getenv("COOKIE_SECURE") == "true",
		CorsOrigin:       origin,
		Port:             port,
	}
}
