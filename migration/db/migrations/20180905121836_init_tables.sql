
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
	id SERIAL PRIMARY KEY, 
	email VARCHAR(100) NOT NULL,
	password VARCHAR(100) NOT NULL,
	is_admin BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_quotas (
	id SERIAL PRIMARY KEY, 
	user_id INTEGER NOT NULL, 
	quota INTEGER NOT NULL, 
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE resources (
	id SERIAL PRIMARY KEY, 
	content VARCHAR(5) NOT NULL, 
	created_by INTEGER NOT NULL, 
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE users;
DROP TABLE user_quotas;
DROP TABLE resources;
