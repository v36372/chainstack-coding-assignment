
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE user_quotas ADD COLUMN current_quota_left INTEGER CONSTRAINT valid_current_quota CHECK(current_quota_left >= 0 AND current_quota_left <= quota);
ALTER TABLE user_quotas ADD CONSTRAINT valid_quota CHECK (quota >= 0);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE user_quotas DROP COLUMN current_quota_left;
