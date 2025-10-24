-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username           VARCHAR(50) UNIQUE  NOT NULL,
    email              VARCHAR(100) UNIQUE NOT NULL,
    password_hash      VARCHAR(255)        NOT NULL,
    role               VARCHAR(20) NOT NULL DEFAULT 'unverified',
    verification_token VARCHAR(100),
    created_at         TIMESTAMPTZ      DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

DROP EXTENSION IF EXISTS "uuid-ossp";


DROP TABLE IF EXISTS users;
-- +goose StatementEnd
