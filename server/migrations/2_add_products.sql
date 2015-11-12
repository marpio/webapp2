-- +migrate Up
CREATE TABLE products (id TEXT PRIMARY KEY NOT NULL, name TEXT NOT NULL, description TEXT NOT NULL, img_path TEXT NOT NULL,  price REAL NOT NULL);
