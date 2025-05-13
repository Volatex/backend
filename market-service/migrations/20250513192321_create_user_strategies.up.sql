CREATE TABLE IF NOT EXISTS user_strategies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    figi TEXT NOT NULL,
    buy_price DOUBLE PRECISION,
    buy_quantity INTEGER,
    sell_price DOUBLE PRECISION,
    sell_quantity INTEGER,
    tinkoff_token TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);