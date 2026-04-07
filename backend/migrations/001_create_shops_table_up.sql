CREATE TABLE IF NOT EXISTS shops (
    id         BIGSERIAL    PRIMARY KEY,
    name       TEXT         NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_shops_deleted_at ON shops (deleted_at);
