// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	"github.com/abdulazizax/ai-embedding/config"
	_ "github.com/abdulazizax/ai-embedding/docs"
	"github.com/abdulazizax/ai-embedding/internal/controller/http/v1/handler"
	"github.com/abdulazizax/ai-embedding/internal/usecase"
	"github.com/abdulazizax/ai-embedding/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description This is a sample server Go Clean Template server.
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase) {
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase)

	// Swagger
	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routes
	v1 := engine.Group("/v1")

	movie := v1.Group("/movie")
	{
		movie.POST("/", handlerV1.CreateMovie)
		movie.GET("/list", handlerV1.GetMovies)
		movie.GET("/:id", handlerV1.GetMovie)
		movie.PUT("/", handlerV1.UpdateMovie)
		movie.DELETE("/:id", handlerV1.DeleteMovie)
		movie.GET("/search", handlerV1.SearchMovie)
	}
}
