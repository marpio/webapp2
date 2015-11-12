-- +migrate Up
CREATE TABLE orders (id TEXT PRIMARY KEY NOT NULL,
                     userId TEXT NOT NULL,
                     state TEXT NOT NULL, 
                     createdAt TEXT NOT NULL,
                     changedAt TEXT NOT NULL,
                     FOREIGN KEY(userId) REFERENCES users(id));
