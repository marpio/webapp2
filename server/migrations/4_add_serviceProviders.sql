-- +migrate Up
CREATE TABLE serviceProviders (id TEXT PRIMARY KEY NOT NULL, email TEXT NOT NULL, phone TEXT NOT NULL, fname TEXT NOT NULL, lname TEXT NOT NULL, postalCode TEXT NOT NULL, address TEXT NOT NULL, description TEXT NOT NULL, img_path TEXT NOT NULL,  price_m2 REAL NOT NULL, token_type TEXT NOT NULL, stripe_publishable_key TEXT NOT NULL, scope TEXT NOT NULL, livemode TEXT NOT NULL, stripe_user_id TEXT NOT NULL, refresh_token TEXT NOT NULL, access_token TEXT NOT NULL);
CREATE UNIQUE INDEX idx_serviceProviders_postalCode on serviceProviders (postalCode);
