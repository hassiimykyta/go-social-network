package main

import (
	"context"
	"fmt"
	"go-rest-chi/internal/config"
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/httpserver"
	"go-rest-chi/internal/router"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func openDB(cfg *config.Config) (*appdb.SQL, error) {
	sqlc, err := appdb.Open(appdb.Options{
		Driver:      cfg.DB.Driver,
		DSN:         cfg.DB.DSN,
		MaxOpen:     cfg.DB.MaxOpen,
		MaxIdle:     cfg.DB.MaxIdle,
		MaxIdleTime: cfg.DB.MaxIdleTime,
	})
	if err != nil {
		return nil, err
	}

	return sqlc, nil
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)

	sqlDB, err := openDB(cfg)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer sqlDB.Close()

	r := router.New(router.Deps{
		DB: sqlDB,
	}, router.Options{
		CORS: router.CORSOpts{
			AllowedOrigins:   cfg.CORS.AllowedOrigins,
			AllowedMethods:   cfg.CORS.AllowedMethods,
			AllowedHeaders:   cfg.CORS.AllowedHeaders,
			ExposedHeaders:   cfg.CORS.ExposedHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           cfg.CORS.MaxAge,
		},
	})

	srv, err := httpserver.New(httpserver.Options{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	})

	if err != nil {
		log.Fatalf("httpserver.New : %v", err)
	}

	log.Printf("▶ listening on http://%s (env=%s)", addr, cfg.App.Env)
	srv.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("⏳ shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("✅ bye")

}
