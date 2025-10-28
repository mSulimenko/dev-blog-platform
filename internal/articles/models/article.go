package models

import "time"

type Article struct {
	Id        string
	Title     string
	Content   string
	AuthorId  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListArticleParams struct {
	AuthorId *string
	Status   *string
	Offset   int
	Limit    int
}

type CreateArticleParams struct {
	Title    string
	Content  string
	AuthorId string
	Status   string
}

type UpdateArticleParams struct {
	Title   *string
	Content *string
	Status  *string
}
