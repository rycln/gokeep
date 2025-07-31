-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
ADD COLUMN salt VARCHAR(16) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users 
DROP COLUMN salt;
-- +goose StatementEnd