-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username           VARCHAR(50) UNIQUE  NOT NULL,
    email              VARCHAR(100) UNIQUE NOT NULL,
    password_hash      VARCHAR(255)        NOT NULL,
    email_verified     BOOLEAN     DEFAULT false,
    verification_token VARCHAR(100),
    created_at         TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS articles
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title      VARCHAR(255) NOT NULL,
    content    TEXT         NOT NULL,
    author_id  UUID         NOT NULL,
    status     VARCHAR(20) DEFAULT 'draft',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_articles_author
        FOREIGN KEY (author_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_articles_status_created;
DROP INDEX IF EXISTS idx_articles_author_id;

DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd