
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE users ADD COLUMN created_by INTEGER;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE users DROP COLUMN created_by;
