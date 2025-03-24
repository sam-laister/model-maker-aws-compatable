-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks ADD COLUMN metadata JSON;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks DROP COLUMN metadata;
-- +goose StatementEnd