package service

import (
	"context"
	"fmt"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/dto"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/models"
	"go.uber.org/zap"
)

type ArticlesRepo interface {
	CreateArticle(ctx context.Context, params models.CreateArticleParams) (*models.Article, error)
	UpdateArticle(ctx context.Context, id string, update models.UpdateArticleParams) (*models.Article, error)
	GetArticleById(ctx context.Context, id string) (*models.Article, error)
	DeleteArticle(ctx context.Context, id string) error
	ListArticles(ctx context.Context, params models.ListArticleParams) ([]*models.Article, error)
}

type ArticlesService struct {
	log  *zap.SugaredLogger
	repo ArticlesRepo
}

func NewArticlesService(log *zap.SugaredLogger, repo ArticlesRepo) *ArticlesService {
	return &ArticlesService{
		log:  log,
		repo: repo,
	}
}

func (a *ArticlesService) CreateArticle(ctx context.Context, req dto.CreateReq) (*dto.ArticleResp, error) {
	params := models.CreateArticleParams{
		Title:    req.Title,
		Content:  req.Content,
		AuthorId: req.AuthorId,
		Status:   req.Status,
	}

	article, err := a.repo.CreateArticle(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed Create Article: %w", err)
	}
	resp := dto.FromArticleModel(article)
	return &resp, nil
}

func (a *ArticlesService) GetArticle(ctx context.Context, id string) (*dto.ArticleResp, error) {
	article, err := a.repo.GetArticleById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed GetArticleById: %w", err)
	}

	resp := dto.FromArticleModel(article)
	return &resp, nil
}

func (a *ArticlesService) DeleteArticle(ctx context.Context, articleId string, userID string) error {
	article, err := a.repo.GetArticleById(ctx, articleId)
	if err != nil {
		return fmt.Errorf("failed to get article: %w", err)
	}

	if article.AuthorId != userID {
		return fmt.Errorf("forbidden: only author can delete article")
	}

	err = a.repo.DeleteArticle(ctx, articleId)
	if err != nil {
		return fmt.Errorf("failed DeleteArticle: %w", err)
	}
	return nil
}

func (a *ArticlesService) ListArticles(ctx context.Context, req dto.ListReq) (*dto.ListResp, error) {

	params := models.ListArticleParams{
		AuthorId: req.AuthorID,
		Status:   req.Status,
		Offset:   req.Offset,
		Limit:    req.Limit,
	}

	articles, err := a.repo.ListArticles(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed ListArticles: %w", err)
	}

	response := dto.FromArticleModels(articles, req.Offset, req.Limit)
	return &response, nil
}

func (a *ArticlesService) UpdateArticle(ctx context.Context,
	articleId string,
	req dto.UpdateReq,
	userID string,
) (*dto.ArticleResp, error) {
	article, err := a.repo.GetArticleById(ctx, articleId)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	if article.AuthorId != userID {
		return nil, fmt.Errorf("forbidden: only author can update article")
	}

	params := models.UpdateArticleParams{
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
	}

	updatedArticle, err := a.repo.UpdateArticle(ctx, articleId, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	response := dto.FromArticleModel(updatedArticle)
	return &response, nil
}
