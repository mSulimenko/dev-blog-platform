package service

import (
	"context"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/models"
)

type ArticlesRepo interface {
	CreateArticle(ctx context.Context, params models.CreateArticleParams) (*models.Article, error)
	UpdateArticle(ctx context.Context, id string, update models.UpdateArticleParams) (*models.Article, error)
	GetArticle(ctx context.Context, id string) (*models.Article, error)
	DeleteArticle(ctx context.Context, id string) error
	ListArticles(ctx context.Context, params models.ListArticleParams) ([]*models.Article, error)
}
