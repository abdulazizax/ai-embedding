package handler

import (
	"strconv"

	"github.com/abdulazizax/ai-embedding/config"
	"github.com/abdulazizax/ai-embedding/internal/entity"
	"github.com/gin-gonic/gin"
)

// CreateMovie godoc
// @Router /movie [post]
// @Summary Create a new movie
// @Description Create a new movie
// @Security BearerAuth
// @Tags movie
// @Accept  json
// @Produce  json
// @Param movie body entity.Movie true "Movie object"
// @Success 201 {object} entity.Movie
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateMovie(ctx *gin.Context) {
	var (
		body entity.Movie
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	movie, err := h.UseCase.MovieRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating movie") {
		return
	}

	ctx.JSON(201, movie)
}

// GetMovie godoc
// @Router /movie/{id} [get]
// @Summary Get a movie by ID
// @Description Get a movie by ID
// @Security BearerAuth
// @Tags movie
// @Accept  json
// @Produce  json
// @Param id path string true "Movie ID"
// @Success 200 {object} entity.Movie
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetMovie(ctx *gin.Context) {
	var (
		req entity.MovieSingleRequest
	)

	req.ID = ctx.Param("id")

	movie, err := h.UseCase.MovieRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting movie") {
		return
	}

	ctx.JSON(200, movie)
}

// GetMovies godoc
// @Router /movie/list [get]
// @Summary Get a list of users
// @Description Get a list of users
// @Security BearerAuth
// @Tags movie
// @Accept  json
// @Produce  json
// @Param page query number false "page"
// @Param limit query number false "limit"
// @Success 200 {object} entity.MovieList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetMovies(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	movies, err := h.UseCase.MovieRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting movie") {
		return
	}

	ctx.JSON(200, movies)
}

// UpdateMovie godoc
// @Router /movie [put]
// @Summary Update a movie
// @Description Update a movie
// @Security BearerAuth
// @Tags movie
// @Accept  json
// @Produce  json
// @Param movie body entity.Movie true "Movie object"
// @Success 200 {object} entity.Movie
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateMovie(ctx *gin.Context) {
	var (
		body entity.Movie
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	movie, err := h.UseCase.MovieRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating movie") {
		return
	}

	ctx.JSON(200, movie)
}

// DeleteMovie godoc
// @Router /movie/{id} [delete]
// @Summary Delete a movie
// @Description Delete a movie
// @Security BearerAuth
// @Tags movie
// @Accept  json
// @Produce  json
// @Param id path string true "Movie ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteMovie(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.MovieRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting movie") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Movie deleted successfully",
	})
}

// SearchMovie godoc
// @Router /movie/search [get]
// @Summary Get movies by search query
// @Description Get movies by search query
// @Tags movie
// @Accept  json
// @Produce  json
// @Param search query string false "Search query"
// @Success 200 {array} entity.Movie
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) SearchMovie(ctx *gin.Context) {
	var (
		req entity.MovieSingleRequest
	)

	search := ctx.DefaultQuery("search", "")

	if search != "" {
		req.NameUz = search
		req.NameRu = search
		req.NameEn = search
	}

	resp, err := h.UseCase.MovieRepo.Search(ctx, req)
	if h.HandleDbError(ctx, err, "Error searching movie") {
		return
	}

	ctx.JSON(200, resp)
}
