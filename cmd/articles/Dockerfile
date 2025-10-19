FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o articles-service ./cmd/articles

EXPOSE 8082

CMD ["./articles-service"]