CREATE TABLE users (
	id uuid NOT NULL PRIMARY KEY,
	email varchar NOT NULL UNIQUE,
	created_at timestamptz NOT NULL
);
