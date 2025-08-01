-- +goose Up
-- +goose StatementBegin
ALTER TABLE items 
RENAME COLUMN deleted TO is_deleted;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items 
RENAME COLUMN is_deleted TO deleted;
-- +goose StatementEnd
