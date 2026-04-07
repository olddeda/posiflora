CREATE TABLE IF NOT EXISTS telegram_integrations (
    id         BIGSERIAL    PRIMARY KEY,
    shop_id    BIGINT       NOT NULL REFERENCES shops (id) ON DELETE CASCADE,
    bot_token  TEXT         NOT NULL,
    chat_id    TEXT         NOT NULL,
    enabled    BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uni_telegram_integrations_shop_id UNIQUE (shop_id)
);

CREATE INDEX IF NOT EXISTS idx_telegram_integrations_deleted_at ON telegram_integrations (deleted_at);
