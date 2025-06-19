package config

import (
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Session  SessionConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	StaticDir    string
	TemplateDir  string
}

type DatabaseConfig struct {
	Driver string
	Path   string
}

type SessionConfig struct {
	CookieName     string
	SessionTTL     time.Duration
	CookieHTTPOnly bool
	CookieSecure   bool
}

func LoadConfig() Config {
	return Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			StaticDir:    getEnv("STATIC_DIR", "./web/static"),
			TemplateDir:  getEnv("TEMPLATE_DIR", "./web/templates"),
		},
		Database: DatabaseConfig{
			Driver: getEnv("DATABASE_DRIVER", "sqlite3"),
			Path:   getEnv("DB_PATH", "./store.db"),
		},
		Session: SessionConfig{
			CookieName:     getEnv("COOKIE_NAME", "session_id"),
			SessionTTL:     24 * time.Hour,
			CookieHTTPOnly: true,
			CookieSecure:   false,
		},
	}
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
