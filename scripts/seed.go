package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	connStr := "postgres://user:password@localhost:5432/dev_blog?sslmode=disable"

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close(ctx)

	log.Println("Connected to database successfully")

	// Засеиваем данные
	if err := seedUsers(ctx, conn); err != nil {
		log.Fatal("Failed to seed users:", err)
	}

	if err := seedArticles(ctx, conn); err != nil {
		log.Fatal("Failed to seed articles:", err)
	}

	log.Println("Seed data inserted successfully!")
}

func seedUsers(ctx context.Context, conn *pgx.Conn) error {
	users := []struct {
		id       string
		username string
		email    string
	}{
		{"11111111-1111-1111-1111-111111111111", "test_user_1", "user1@test.com"},
		{"22222222-2222-2222-2222-222222222222", "test_user_2", "user2@test.com"},
		{"33333333-3333-3333-3333-333333333333", "test_user_3", "user3@test.com"},
	}

	for _, u := range users {
		_, err := conn.Exec(ctx, `
			INSERT INTO users (id, username, email, password_hash, role) 
			VALUES ($1, $2, $3, $4, 'user')
			ON CONFLICT (id) DO NOTHING`,
			u.id, u.username, u.email, "fake_password_hash_"+u.id,
		)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", u.username, err)
		}
		log.Printf("User seeded: %s (%s)", u.username, u.email)
	}
	return nil
}

func seedArticles(ctx context.Context, conn *pgx.Conn) error {
	articles := []struct {
		id        string
		title     string
		content   string
		authorID  string
		status    string
		createdAt time.Time
	}{
		// Статьи пользователя 1 (самые старые)
		{"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1", "Первая статья", "Содержание первой статьи...", "11111111-1111-1111-1111-111111111111", "published", time.Now().Add(-2 * time.Hour)},
		{"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2", "Вторая статья", "Содержание второй статьи...", "11111111-1111-1111-1111-111111111111", "published", time.Now().Add(-90 * time.Minute)},
		{"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3", "Третья статья", "Содержание третьей статьи...", "11111111-1111-1111-1111-111111111111", "published", time.Now().Add(-80 * time.Minute)},

		// Статьи пользователя 2
		{"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1", "Статья пользователя 2", "Отличный контент от второго пользователя...", "22222222-2222-2222-2222-222222222222", "published", time.Now().Add(-70 * time.Minute)},
		{"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2", "Еще одна статья", "Интересные мысли и идеи...", "22222222-2222-2222-2222-222222222222", "published", time.Now().Add(-60 * time.Minute)},

		// Статьи пользователя 3
		{"cccccccc-cccc-cccc-cccc-ccccccccccc1", "Новейшая статья", "Самый свежий контент на платформе!", "33333333-3333-3333-3333-333333333333", "published", time.Now().Add(-50 * time.Minute)},
		{"cccccccc-cccc-cccc-cccc-ccccccccccc2", "Техническая статья", "Глубокий разбор современных технологий...", "33333333-3333-3333-3333-333333333333", "published", time.Now().Add(-40 * time.Minute)},

		// Дополнительные статьи чтобы было больше 10
		{"dddddddd-dddd-dddd-dddd-ddddddddddd1", "Статья 8", "Интересный контент восьмой статьи...", "11111111-1111-1111-1111-111111111111", "published", time.Now().Add(-30 * time.Minute)},
		{"dddddddd-dddd-dddd-dddd-ddddddddddd2", "Статья 9", "Полезная информация девятой статьи...", "22222222-2222-2222-2222-222222222222", "published", time.Now().Add(-20 * time.Minute)},
		{"dddddddd-dddd-dddd-dddd-ddddddddddd3", "Статья 10", "Завершающая десятая статья...", "33333333-3333-3333-3333-333333333333", "published", time.Now().Add(-10 * time.Minute)},
		{"dddddddd-dddd-dddd-dddd-ddddddddddd4", "Статья 11", "Лишняя статья для проверки лимита...", "11111111-1111-1111-1111-111111111111", "published", time.Now().Add(-5 * time.Minute)},
	}

	for _, a := range articles {
		_, err := conn.Exec(ctx, `
			INSERT INTO articles (id, title, content, author_id, status, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (id) DO NOTHING`,
			a.id, a.title, a.content, a.authorID, a.status, a.createdAt, a.createdAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert article %s: %w", a.title, err)
		}
		log.Printf("Article seeded: %s (author: %s)", a.title, a.authorID)
	}

	// Проверяем сколько статей создано
	var count int
	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM articles").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count articles: %w", err)
	}

	log.Printf("Total articles in database: %d", count)
	return nil
}
