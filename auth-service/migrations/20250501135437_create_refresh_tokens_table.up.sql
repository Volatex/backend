CREATE TABLE refresh_tokens (
    token VARCHAR(255) PRIMARY KEY,
    user_id BIGSERIAL REFERENCES users(id),
    expires_at TIMESTAMPTZ NOT NULL
)