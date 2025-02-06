package usecase

import (
	"github.com/abdulazizax/ai-embedding/config"
	"github.com/abdulazizax/ai-embedding/internal/usecase/repo"
	"github.com/abdulazizax/ai-embedding/pkg/logger"
	"github.com/abdulazizax/ai-embedding/pkg/postgres"
	openai "github.com/sashabaranov/go-openai"
)

// UseCase -.
type UseCase struct {
	MovieRepo MovieRepoI
}

// New -.
func New(openaiClient *openai.Client, pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		MovieRepo: repo.NewMovieRepo(openaiClient, pg, config, logger),
	}
}
