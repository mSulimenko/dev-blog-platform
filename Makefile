run-articles:
	go run cmd/articles/main.go

run-auth:
	go run cmd/auth/main.go

run-notify:
	go run cmd/notify/main.go

goose-down:
	goose -dir migrations postgres "postgres://user:password@localhost:5432/dev_blog?sslmode=disable" down

gen-proto:
	protoc -I protos protos/auth.proto \
        --go_out=./protos/gen/go/ \
        --go_opt=paths=source_relative \
        --go-grpc_out=./protos/gen/go/ \
        --go-grpc_opt=paths=source_relative