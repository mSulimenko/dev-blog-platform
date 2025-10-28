package dto

import (
	"github.com/mSulimenko/dev-blog-platform/internal/articles/models"
	"time"
)

type CreateReq struct {
	Title    string `json:"title" validate:"required,min=1,max=255"`
	Content  string `json:"content" validate:"required,min=1"`
	Status   string `json:"status" validate:"required,oneof=draft published archived"`
	AuthorId string `json:"author_id" validate:"required"`
}

type ListReq struct {
	AuthorID *string `json:"author_id" validate:"omitempty"`
	Status   *string `json:"status" validate:"omitempty,oneof=draft published archived"`
	Offset   int     `json:"offset" validate:"min=0"`
	Limit    int     `json:"limit" validate:"min=1,max=100"`
}

type ListResp struct {
	Articles []ArticleResp `json:"articles"`
	Total    int           `json:"total"`
	Offset   int           `json:"offset"`
	Limit    int           `json:"limit"`
}

type UpdateReq struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content *string `json:"content,omitempty" validate:"omitempty,min=1"`
	Status  *string `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

type ArticleResp struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromArticleModel(article *models.Article) ArticleResp {
	return ArticleResp{
		ID:        article.Id,
		Title:     article.Title,
		Content:   article.Content,
		AuthorID:  article.AuthorId,
		Status:    article.Status,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}
}

func FromArticleModels(articles []*models.Article, offset, limit int) ListResp {
	articleResponses := make([]ArticleResp, len(articles))
	for i, article := range articles {
		articleResponses[i] = FromArticleModel(article)
	}

	return ListResp{
		Articles: articleResponses,
		Total:    len(articles),
		Offset:   offset,
		Limit:    limit,
	}
}
