package httphandler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/dto"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/transport/grpc"
	"go.uber.org/zap"
)

type ArticlesServiceInterface interface {
	CreateArticle(ctx context.Context, req dto.CreateRequest) (*dto.ArticleResponse, error)
	GetArticle(ctx context.Context, id string) (*dto.ArticleResponse, error)
	DeleteArticle(ctx context.Context, articleId string, userID string) error
	ListArticles(ctx context.Context, req dto.ListRequest) (*dto.ListResponse, error)
	UpdateArticle(ctx context.Context, articleId string, req dto.UpdateRequest, userID string) (*dto.ArticleResponse, error)
}

type Handler struct {
	articlesService ArticlesServiceInterface
	log             *zap.SugaredLogger
	validate        *validator.Validate
	authClient      *grpc.Client
}

func NewHandler(articlesService ArticlesServiceInterface, logger *zap.SugaredLogger, authClient *grpc.Client) *Handler {
	validate := validator.New()
	return &Handler{
		articlesService: articlesService,
		log:             logger,
		validate:        validate,
		authClient:      authClient,
	}
}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/articles", func(r chi.Router) {
			// Публичные endpoints
			r.Get("/", h.ListArticles)   // GET /api/v1/articles
			r.Get("/{id}", h.GetArticle) // GET /api/v1/articles/{id}

			// Защищенные endpoints
			r.Group(func(r chi.Router) {
				r.Use(h.AuthMiddleware(h.authClient))

				r.Post("/", h.CreateArticle)       // POST /api/v1/articles
				r.Put("/{id}", h.UpdateArticle)    // PUT /api/v1/articles/{id}
				r.Delete("/{id}", h.DeleteArticle) // DELETE /api/v1/articles/{id}
			})
		})
	})

	return router
}
