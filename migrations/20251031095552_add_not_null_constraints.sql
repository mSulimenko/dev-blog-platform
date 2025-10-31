-- +goose Up
-- +goose StatementBegin

-- 1. Заполняем существующие NULL значения
UPDATE users SET verification_token = '' WHERE verification_token IS NULL;

-- 2. Добавляем NOT NULL constraint
ALTER TABLE users
    ALTER COLUMN verification_token SET NOT NULL,
    ALTER COLUMN verification_token SET DEFAULT '';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Убираем NOT NULL и DEFAULT
ALTER TABLE users
    ALTER COLUMN verification_token DROP NOT NULL,
    ALTER COLUMN verification_token DROP DEFAULT;

-- +goose StatementEnd