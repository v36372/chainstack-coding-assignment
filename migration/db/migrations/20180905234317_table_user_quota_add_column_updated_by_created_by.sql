
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE user_quotas ADD COLUMN created_by INTEGER;
ALTER TABLE user_quotas ADD COLUMN updated_by INTEGER;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE user_quotas DROP COLUMN created_by;
ALTER TABLE user_quotas DROP COLUMN updated_by;
