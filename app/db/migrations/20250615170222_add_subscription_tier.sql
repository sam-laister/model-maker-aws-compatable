-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN subscription_tier TEXT NOT NULL DEFAULT 'free';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN subscription_tier;
-- +goose StatementEnd 