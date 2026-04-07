CREATE TABLE IF NOT EXISTS telegram_send_logs (
    id         BIGSERIAL    PRIMARY KEY,
    shop_id    BIGINT       NOT NULL REFERENCES shops (id)  ON DELETE CASCADE,
    order_id   BIGINT       NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    message    TEXT         NOT NULL,
    status     VARCHAR(10)  NOT NULL,
    error      TEXT,
    sent_at    TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT idx_shop_order UNIQUE (shop_id, order_id)
);

CREATE INDEX IF NOT EXISTS idx_telegram_send_logs_deleted_at ON telegram_send_logs (deleted_at);
