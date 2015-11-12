-- +migrate Up
CREATE TABLE productOrders (id TEXT PRIMARY KEY NOT NULL,
                            orderId TEXT NOT NULL,
                            productId TEXT NOT NULL, 
                            ammount REAL NOT NULL,
                            FOREIGN KEY(productId) REFERENCES products(id),
                            FOREIGN KEY(orderId) REFERENCES orders(id));
