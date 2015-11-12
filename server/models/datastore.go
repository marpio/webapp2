package models

import (
	"github.com/jmoiron/sqlx"
	vm "github.com/marpio/webapp2/server/viewmodels"
)

type Datastore interface {
	AllProducts() ([]*ProductRow, error)
	GetProductByID(id string) (*ProductRow, error)
	CreateOrder(order *vm.OrderDto) (*vm.OrderDto, error)
	GetOrderByID(id string) (*vm.OrderDto, error)
	AllUsers() ([]*UserRow, error)
	GetUserByEmail(email string) (*UserRow, error)
	GetUserByEmailAndPassword(email, password string) (*UserRow, error)
	Signup(email, password, passwordAgain string) (*UserRow, error)
	GetUserByID(id string) (*UserRow, error)
}

type datastore struct {
	db *sqlx.DB
}

func NewDatastore(db *sqlx.DB) Datastore {
	ds := &datastore{
		db: db,
	}
	return ds
}
