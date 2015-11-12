package models

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const usersTablename string = "users"

type userTable struct {
	Base
}

func newUserTable(ds *datastore) *userTable {
	u := &userTable{}
	u.table = usersTablename
	u.db = ds.db
	return u
}

type UserRow struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (ds *datastore) AllUsers() ([]*UserRow, error) {
	users := []*UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v", usersTablename)
	err := ds.db.Select(&users, query)

	return users, err
}

func (ds *datastore) GetUserByID(id string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", usersTablename)
	err := ds.db.Get(user, query, id)

	return user, err
}

func (ds *datastore) GetUserByEmail(email string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE email=$1", usersTablename)
	err := ds.db.Get(user, query, email)

	return user, err
}

// GetByEmail returns record by email but checks password first.
func (ds *datastore) GetUserByEmailAndPassword(email, password string) (*UserRow, error) {
	user, err := ds.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, err
}

// Signup create a new record of user.
func (ds *datastore) Signup(email, password, passwordAgain string) (*UserRow, error) {
	if email == "" {
		return nil, errors.New("Email cannot be blank.")
	}
	if password == "" {
		return nil, errors.New("Password cannot be blank.")
	}
	if password != passwordAgain {
		return nil, errors.New("Password is invalid.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["id"] = uuid.NewV4().String()
	data["email"] = email
	data["password"] = hashedPassword
	u := newUserTable(ds)
	sqlResult, err := u.InsertIntoTable(nil, data)
	if err != nil {
		return nil, err
	}
	userId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return ds.GetUserByID(userId)
}
