-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
ALTER COLUMN salt TYPE VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users 
ALTER COLUMN salt TYPE VARCHAR(16);
-- +goose StatementEnd