CREATE TABLE IF NOT EXISTS user_strategies (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_tokens(user_id),
    figi TEXT NOT NULL,
    buy_price DOUBLE PRECISION,
    buy_quantity INTEGER,
    sell_price DOUBLE PRECISION,
    sell_quantity INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
