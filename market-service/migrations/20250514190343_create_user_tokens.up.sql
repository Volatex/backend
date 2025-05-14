CREATE TABLE IF NOT EXISTS user_tokens (
    user_id UUID PRIMARY KEY,
    tinkoff_token TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
