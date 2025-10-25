package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/models"
	"strings"
	"time"
)

type ArticlesRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *ArticlesRepository {
	return &ArticlesRepository{
		db: db,
	}
}

func (a *ArticlesRepository) CreateArticle(ctx context.Context, params models.CreateArticleParams) (*models.Article, error) {
	article := models.Article{
		Title:    params.Title,
		Content:  params.Content,
		AuthorId: params.AuthorId,
		Status:   params.Status,
	}
	q := `INSERT INTO articles(title, content, author_id, status) 
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`

	err := a.db.QueryRow(ctx, q, params.Title, params.Content, params.AuthorId, params.Status).
		Scan(&article.Id, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed QueryRow: %w", err)
	}

	return &article, nil
}

func (a *ArticlesRepository) UpdateArticle(ctx context.Context, id string, update models.UpdateArticleParams) (*models.Article, error) {
	var article models.Article
	q, args := a.buildUpdateQuery(id, update)

	err := a.db.QueryRow(ctx, q, args...).Scan(
		&article.Id,
		&article.Title,
		&article.Content,
		&article.AuthorId,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrArticleNotFound
		}
		return nil, err
	}

	return &article, nil
}

func (a *ArticlesRepository) buildUpdateQuery(id string, updates models.UpdateArticleParams) (string, []interface{}) {
	query := `UPDATE articles SET updated_at = $1`
	args := []interface{}{time.Now()}

	if updates.Title != nil {
		args = append(args, *updates.Title)
		query += fmt.Sprintf(", title = $%d", len(args))
	}
	if updates.Content != nil {
		args = append(args, *updates.Content)
		query += fmt.Sprintf(", content = $%d", len(args))
	}
	if updates.Status != nil {
		args = append(args, *updates.Status)
		query += fmt.Sprintf(", status = $%d", len(args))
	}

	args = append(args, id)
	query += fmt.Sprintf(
		` WHERE id = $%d RETURNING id, title, content, author_id, status, created_at, updated_at`, len(args))

	return query, args
}

func (a *ArticlesRepository) GetArticle(ctx context.Context, id string) (*models.Article, error) {
	var article models.Article
	q := `SELECT id, title, content, author_id, status, created_at, updated_at
			FROM articles
			WHERE id = $1`

	err := a.db.QueryRow(ctx, q, id).Scan(
		&article.Id,
		&article.Title,
		&article.Content,
		&article.AuthorId,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrArticleNotFound
		}
		return nil, err
	}
	return &article, nil
}

func (a *ArticlesRepository) DeleteArticle(ctx context.Context, id string) error {
	q := `DELETE FROM articles WHERE id = $1`
	_, err := a.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticlesRepository) ListArticles(ctx context.Context, params models.ListArticleParams) ([]*models.Article, error) {
	var articles []*models.Article

	q, values := a.buildListQuery(params)

	rows, err := a.db.Query(ctx, q, values...)
	if err != nil {
		return nil, fmt.Errorf("querying articles: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var art models.Article
		err = rows.Scan(
			&art.Id,
			&art.Title,
			&art.Content,
			&art.AuthorId,
			&art.Status,
			&art.CreatedAt,
			&art.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning article: %w", err)
		}
		articles = append(articles, &art)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return articles, nil

}

func (a *ArticlesRepository) buildListQuery(reqParams models.ListArticleParams) (string, []interface{}) {
	var whereClauses []string
	var values []interface{}

	if reqParams.AuthorId != nil {
		values = append(values, *reqParams.AuthorId)
		whereClauses = append(whereClauses, fmt.Sprintf("author_id = $%d", len(values)))
	}

	if reqParams.Status != nil {
		values = append(values, *reqParams.Status)
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", len(values)))
	}

	q := `SELECT id, title, content, author_id, status, created_at, updated_at FROM articles`

	if len(whereClauses) > 0 {
		q += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	q += " ORDER BY created_at DESC"

	values = append(values, reqParams.Limit, reqParams.Offset)
	q += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(values)-1, len(values))

	return q, values
}
