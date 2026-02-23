-- +goose Up
CREATE TABLE orders (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    loan_id     VARCHAR(64) NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'created',
    total_amount BIGINT NOT NULL,
    currency    VARCHAR(3) NOT NULL DEFAULT 'SAR',
    card_token  VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);

-- +goose Down
DROP TABLE IF EXISTS orders;
