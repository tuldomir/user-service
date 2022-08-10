CREATE TABLE users (
	id varchar NOT NULL PRIMARY KEY,
	email varchar NOT NULL UNIQUE,
	created_at timestamptz NOT NULL DEFAULT (now())
);