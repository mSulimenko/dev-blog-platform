run-articles:
	go run cmd/articles/main.go

goose-down:
	goose -dir migrations postgres "postgres://user:password@localhost:5432/dev_blog?sslmode=disable" down
