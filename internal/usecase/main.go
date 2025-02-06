package usecase

import (
	"github.com/abdulazizax/ai-embedding/config"
	"github.com/abdulazizax/ai-embedding/internal/usecase/repo"
	"github.com/abdulazizax/ai-embedding/pkg/logger"
	"github.com/abdulazizax/ai-embedding/pkg/postgres"
	"github.com/google/generative-ai-go/genai"
)

// UseCase -.
type UseCase struct {
	MovieRepo MovieRepoI
}

// New -.
func New(genaiClient *genai.Client, pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		MovieRepo: repo.NewMovieRepo(genaiClient, pg, config, logger),
	}
}
