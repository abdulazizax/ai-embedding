// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/abdulazizax/ai-embedding/config"
	v1 "github.com/abdulazizax/ai-embedding/internal/controller/http/v1"
	"github.com/abdulazizax/ai-embedding/internal/usecase"
	"github.com/abdulazizax/ai-embedding/pkg/httpserver"
	"github.com/abdulazizax/ai-embedding/pkg/logger"
	"github.com/abdulazizax/ai-embedding/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	genaiClient, err := genai.NewClient(context.Background(), option.WithAPIKey(cfg.ApiKey))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - genai.NewClient: %w", err))
	}
	defer genaiClient.Close()

	// Use case
	useCase := usecase.New(genaiClient, pg, cfg, l)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, cfg, useCase)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	l.Info("app - Run - httpServer: %s", cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
