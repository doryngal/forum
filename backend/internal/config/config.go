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
	StaticPath   string
	TemplatePath string
}

type DatabaseConfig struct {
	Driver   string
	Path     string
	Provider string
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
			StaticPath:   getEnv("STATIC_PATH", "./../frontend/static"),
			TemplatePath: getEnv("TEMPLATE_PATH", "./../frontend/templates"),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DATABASE_DRIVER", "sqlite3"),
			Path:     getEnv("DB_PATH", "./../backend/forum.db"),
			Provider: getEnv("PROVIDER", "sqlite"),
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
