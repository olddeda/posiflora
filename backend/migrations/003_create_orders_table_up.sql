CREATE TABLE IF NOT EXISTS orders (
    id            BIGSERIAL      PRIMARY KEY,
    shop_id       BIGINT         NOT NULL REFERENCES shops (id) ON DELETE CASCADE,
    number        TEXT           NOT NULL,
    total         NUMERIC(12,2)  NOT NULL,
    customer_name TEXT           NOT NULL,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_orders_shop_id   ON orders (shop_id);
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders (deleted_at);
