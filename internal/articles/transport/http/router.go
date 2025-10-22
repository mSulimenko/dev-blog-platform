package httphandler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/dto"
	"go.uber.org/zap"
)

type UsersServiceInterface interface {
	CreateUser(ctx context.Context, userReq *dto.UserReq) (string, error)
	GetUser(ctx context.Context, id string) (*dto.UserResp, error)
	ListUsers(ctx context.Context) ([]*dto.UserResp, error)
	UpdateUser(ctx context.Context, id string, userReq *dto.UserReq) error
	DeleteUser(ctx context.Context, id string) error
}

type Handler struct {
	usersService UsersServiceInterface
	log          *zap.SugaredLogger
}

func NewHandler(usersService UsersServiceInterface, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		usersService: usersService,
		log:          logger,
	}
}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", h.CreateUser) // POST /api/v1/users
			r.Get("/", h.ListUsers)   // GET /api/v1/users
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetUser)       // GET /api/v1/users/{id}
				r.Put("/", h.UpdateUser)    // PUT /api/v1/users/{id}
				r.Delete("/", h.DeleteUser) // DELETE /api/v1/users/{id}
			})
		})
	})

	return router
}
