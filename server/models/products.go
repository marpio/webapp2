package models

import (
	"fmt"
)

const productsTablename string = "products"

type ProductRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	ImagePath   string `db:"img_path"`
	Price       string `db:"price"`
}

func (ds *datastore) AllProducts() ([]*ProductRow, error) {
	products := []*ProductRow{}
	query := fmt.Sprintf("SELECT * FROM %v", productsTablename)
	err := ds.db.Select(&products, query)

	return products, err
}

func (ds *datastore) GetProductByID(id string) (*ProductRow, error) {
	product := &ProductRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", productsTablename)
	err := ds.db.Get(product, query, id)

	return product, err
}
