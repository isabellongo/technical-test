package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/driva/api/internal/config"
	"github.com/driva/api/internal/handler"
	"github.com/driva/api/internal/middleware"
	"github.com/driva/api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	// Pool de conexões Postgres
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	log.Println("database connected")

	// Repositories
	seedRepo := repository.NewSeedRepository(pool)
	goldRepo := repository.NewGoldRepository(pool)

	// Handlers
	peopleHandler := handler.NewPeopleHandler(seedRepo)
	analyticsHandler := handler.NewAnalyticsHandler(goldRepo)

	// Router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Grupo protegido por API Key
	auth := r.Group("/")
	auth.Use(middleware.APIKeyAuth(cfg.APIKey))
	{
		// Fonte (com simulação 429)
		v1 := auth.Group("/people/v1")
		v1.Use(middleware.Simulate429())
		v1.GET("/enrichments", peopleHandler.GetEnrichments)

		// Analytics (sem 429)
		analytics := auth.Group("/analytics")
		analytics.GET("/overview", analyticsHandler.GetOverview)
		analytics.GET("/enrichments", analyticsHandler.GetEnrichments)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.APIPort,
		Handler: r,
	}

	go func() {
		log.Printf("api listening on :%s", cfg.APIPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Println("api stopped")
}
