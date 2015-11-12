-- +migrate Up
CREATE TABLE users (id TEXT PRIMARY KEY NOT NULL, email TEXT NOT NULL, password TEXT NOT NULL);

CREATE UNIQUE INDEX idx_users_email on users (email);
