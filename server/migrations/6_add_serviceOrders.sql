-- +migrate Up
CREATE TABLE serviceOrders (id TEXT PRIMARY KEY NOT NULL,
                            orderId TEXT NOT NULL,
                            serviceProviderId TEXT NOT NULL,
                            serviceDate TEXT,
                            FOREIGN KEY(orderId) REFERENCES orders(id),
                            FOREIGN KEY(serviceProviderId) REFERENCES serviceProviders(id));
