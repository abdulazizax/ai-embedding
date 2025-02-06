// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/abdulazizax/ai-embedding/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// MovieRepo -.
	MovieRepoI interface {
		Create(ctx context.Context, req entity.Movie) (entity.Movie, error)
		GetSingle(ctx context.Context, req entity.MovieSingleRequest) (entity.Movie, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.MovieList, error)
		Update(ctx context.Context, req entity.Movie) (entity.Movie, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
		Search(ctx context.Context, req entity.MovieSingleRequest) (entity.MovieList, error)
	}
)
