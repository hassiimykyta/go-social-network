package config

import (
	"errors"
	"fmt"
	"go-rest-chi/internal/helpers"
	"log"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env             string
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Driver      string
	DSN         string
	MaxOpen     int
	MaxIdle     int
	MaxIdleTime time.Duration
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type JWTConfig struct {
	Secret     string
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Config struct {
	App  AppConfig
	DB   DBConfig
	CORS CORSConfig
	JWT  JWTConfig
}

func (c *Config) Validate() error {
	if c.App.Env == "prod" {
		if len(c.CORS.AllowedOrigins) == 1 && c.CORS.AllowedOrigins[0] == "*" {
			if c.CORS.AllowCredentials {
				return fmt.Errorf("in prod, CORS_ALLOW_CREDENTIALS=true cannot be used with CORS_ALLOWED_ORIGINS=*")
			}
			log.Println("⚠️  WARNING: in prod with CORS_ALLOWED_ORIGINS=* (no credentials). Consider whitelisting domains.")

		}
	}

	if c.DB.DSN == "" {
		return errors.New("DB_DSN is required")
	}

	switch c.DB.Driver {
	case "postgres", "mysql", "pgx":
	default:
		return fmt.Errorf("unsupported DB_DRIVER %v", c.DB.Driver)
	}

	return nil
}

func (c *Config) RedactedString() string {
	return fmt.Sprintf(
		"env=%s host=%s port=%s db.driver=%s db.maxOpen=%d db.maxIdle=%d",
		c.App.Env, c.App.Host, c.App.Port, c.DB.Driver, c.DB.MaxOpen, c.DB.MaxIdle,
	)
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Env:             helpers.GetEnv("APP_ENV", "dev"),
			Host:            helpers.GetEnv("APP_HOST", "0.0.0.0"),
			Port:            helpers.GetEnv("APP_PORT", "8080"),
			ReadTimeout:     helpers.MustDur(helpers.GetEnv("READ_TIMEOUT", "10s"), 10*time.Second),
			WriteTimeout:    helpers.MustDur(helpers.GetEnv("WRITE_TIMEOUT", "10s"), 10*time.Second),
			IdleTimeout:     helpers.MustDur(helpers.GetEnv("IDLE_TIMEOUT", "60s"), 60*time.Second),
			ShutdownTimeout: helpers.MustDur(helpers.GetEnv("SHUTDOWN_TIMEOUT", "10s"), 10*time.Second),
		},
		DB: DBConfig{
			Driver:      helpers.GetEnv("DB_DRIVER", "postgres"),
			DSN:         helpers.GetEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/yourdb?sslmode=disable"),
			MaxOpen:     helpers.MustInt(helpers.GetEnv("DB_MAX_OPEN", "20"), 20),
			MaxIdle:     helpers.MustInt(helpers.GetEnv("DB_MAX_IDLE", "10"), 10),
			MaxIdleTime: helpers.MustDur(helpers.GetEnv("DB_MAX_IDLE_TIME", "5m"), 5*time.Minute),
		},
		CORS: CORSConfig{
			AllowedOrigins:   helpers.Csv(helpers.GetEnv("CORS_ALLOWED_ORIGINS", "*")),
			AllowedMethods:   helpers.Csv(helpers.GetEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")),
			AllowedHeaders:   helpers.Csv(helpers.GetEnv("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token")),
			ExposedHeaders:   helpers.Csv(helpers.GetEnv("CORS_EXPOSE_HEADERS", "")),
			AllowCredentials: helpers.MustBool(helpers.GetEnv("CORS_ALLOW_CREDENTIALS", "true"), true),
			MaxAge:           helpers.MustInt(helpers.GetEnv("CORS_MAX_AGE", "300"), 300),
		},
		JWT: JWTConfig{
			Secret:     helpers.GetEnv("JWT_SECRET", "dev_secret_change_me"),
			Issuer:     helpers.GetEnv("JWT_ISSUER", "app"),
			AccessTTL:  helpers.MustDur(helpers.GetEnv("JWT_ACCESS_TTL", "15m"), 15*time.Minute),
			RefreshTTL: helpers.MustDur(helpers.GetEnv("JWT_REFRESH_TTL", "720h"), 30*24*time.Hour),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
