CREATE TABLE IF NOT EXISTS users
(
    id                 UUID PRIMARY KEY,
    username           VARCHAR(50) UNIQUE  NOT NULL,
    email              VARCHAR(100) UNIQUE NOT NULL,
    password_hash      VARCHAR(255)        NOT NULL,
    email_verified     BOOLEAN     DEFAULT false,
    verification_token VARCHAR(100),
    created_at         TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS articles
(
    id         UUID PRIMARY KEY,
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

CREATE INDEX IF NOT EXISTS idx_articles_author_id ON articles (author_id);
CREATE INDEX IF NOT EXISTS idx_articles_status_created ON articles (status, created_at DESC);