
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE resources ALTER COLUMN content TYPE VARCHAR(100);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE resources ALTER COLUMN content TYPE VARCHAR(5);

