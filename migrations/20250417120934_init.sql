-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY,
    balance BIGINT NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
