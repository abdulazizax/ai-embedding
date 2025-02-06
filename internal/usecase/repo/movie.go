package repo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abdulazizax/ai-embedding/config"
	"github.com/abdulazizax/ai-embedding/internal/entity"
	"github.com/abdulazizax/ai-embedding/pkg/logger"
	"github.com/abdulazizax/ai-embedding/pkg/postgres"
	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
)

type MovieRepo struct {
	openaiClient *openai.Client
	pg           *postgres.Postgres
	config       *config.Config
	logger       *logger.Logger
}

// New -.
func NewMovieRepo(openaiClient *openai.Client, pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *MovieRepo {
	return &MovieRepo{
		openaiClient: openaiClient,
		pg:           pg,
		config:       config,
		logger:       logger,
	}
}

func (r *MovieRepo) Create(ctx context.Context, req entity.Movie) (entity.Movie, error) {
	req.ID = uuid.NewString()

	embedding, err := r.generateVector(&req)
	if err != nil {
		return entity.Movie{}, err
	}

	formattedEmbedding := formatVectorLiteral(embedding)

	qeury, args, err := r.pg.Builder.Insert("movies").
		Columns(`id, name_uz, name_en, name_ru, embedding`).
		Values(req.ID, req.NameUz, req.NameEn, req.NameRu, formattedEmbedding).ToSql()
	if err != nil {
		return entity.Movie{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Movie{}, err
	}

	return req, nil
}

func (r *MovieRepo) GetSingle(ctx context.Context, req entity.MovieSingleRequest) (entity.Movie, error) {
	response := entity.Movie{}
	var (
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name_uz, name_uz, name_uz, created_at, updated_at`).
		From("movies")

	filters := squirrel.And{}

	if req.ID != "" {
		filters = append(filters, squirrel.Eq{"id": req.ID})
	}
	if req.NameUz != "" {
		filters = append(filters, squirrel.ILike{"name_uz": req.NameUz})
	}
	if req.NameRu != "" {
		filters = append(filters, squirrel.ILike{"name_ru": req.NameRu})
	}
	if req.NameEn != "" {
		filters = append(filters, squirrel.ILike{"name_en": req.NameEn})
	}

	if len(filters) == 0 {
		return entity.Movie{}, fmt.Errorf("GetSingle - invalid request")
	}

	// Add filters to query builder
	qeuryBuilder = qeuryBuilder.Where(filters)

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Movie{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, &response.NameUz, &response.NameEn, &response.NameRu, &createdAt, &updatedAt)
	if err != nil {
		return entity.Movie{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *MovieRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.MovieList, error) {
	var (
		response             = entity.MovieList{}
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name_uz, name_en, name_ru, created_at, updated_at`).
		From("movies")

	qeuryBuilder, where := PrepareGetListQuery(qeuryBuilder, req)

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, qeury, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Movie
		err = rows.Scan(&item.ID, &item.NameUz, &item.NameEn, &item.NameRu, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("movies").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *MovieRepo) Update(ctx context.Context, req entity.Movie) (entity.Movie, error) {
	mp := map[string]interface{}{
		"name_uz":    req.NameUz,
		"name_en":    req.NameEn,
		"name_ru":    req.NameRu,
		"updated_at": "now()",
	}

	qeury, args, err := r.pg.Builder.Update("movies").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Movie{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Movie{}, err
	}

	return req, nil
}

func (r *MovieRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("movies").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *MovieRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("movies").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}

func (r *MovieRepo) Search(ctx context.Context, req entity.MovieSingleRequest) (entity.MovieList, error) {
	var (
		response             = entity.MovieList{}
		createdAt, updatedAt time.Time
		count                int
	)

	embedding, err := r.generateVector(&entity.Movie{
		NameUz: req.NameUz,
		NameEn: req.NameEn,
		NameRu: req.NameRu,
	})
	if err != nil {
		return entity.MovieList{}, err
	}

	formattedEmbedding := formatVectorLiteral(embedding)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name_uz, name_en, name_ru, created_at, updated_at, embedding <-> ? AS distance`).
		From("movies").
		OrderBy("distance").
		Limit(10)

	qeury, _, err := qeuryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, qeury, formattedEmbedding)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Movie
		err = rows.Scan(&item.ID, &item.NameUz, &item.NameEn, &item.NameRu, &createdAt, &updatedAt, &item.Distance)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
		count++
	}

	response.Count = count

	return response, nil
}

func (r *MovieRepo) generateVector(movie *entity.Movie) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{movie.NameUz, movie.NameEn, movie.NameRu},
		Model: openai.AdaEmbeddingV2,
	}
	resp, err := r.openaiClient.CreateEmbeddings(context.Background(), req)
	if err != nil {
		log.Fatalf("Embedding generatsiya qilishda xato: %v", err)
	}

	return resp.Data[0].Embedding, nil
}

func formatVectorLiteral(vector []float32) string {
	// Check if the vector is empty
	if len(vector) == 0 {
		return "{}"
	}

	// Convert the vector values to string
	strValues := make([]string, len(vector))
	for i, v := range vector {
		strValues[i] = fmt.Sprintf("%f", v) // Format float32 to string
	}

	// Join the values with commas and wrap them in square brackets for array
	return "[" + strings.Join(strValues, ",") + "]"
}
