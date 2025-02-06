package handler

import (
	"github.com/abdulazizax/ai-embedding/config"
	"github.com/abdulazizax/ai-embedding/internal/usecase"
	"github.com/abdulazizax/ai-embedding/pkg/logger"
)

type Handler struct {
	Logger  *logger.Logger
	Config  *config.Config
	UseCase *usecase.UseCase
}

func NewHandler(l *logger.Logger, c *config.Config, useCase *usecase.UseCase) *Handler {
	return &Handler{
		Logger:  l,
		Config:  c,
		UseCase: useCase,
	}
}
