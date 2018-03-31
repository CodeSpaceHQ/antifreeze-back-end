package user

import (
	"fmt"

	"cloud.google.com/go/datastore"
	"golang.org/x/crypto/bcrypt"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

type User struct {
	// User Entity's Datastore key
	Key      *datastore.Key `datastore:"__key__"`
	Email    string         `datastore:""`
	Password string         `datastore:"noindex"`
	Devices  []int          `datastore:"noindex"`
}

// In case we want to mock the model for unit tests
type Interface interface {
	GetByEmail(string) (*User, error)
	Put(string, string) error
}

type Model struct {
	*env.Env
}

func (m *Model) GetByEmail(email string) (User, error) {
	// results := make([]User)
	// q := datastore.NewQuery("User")

	// for t := m.Run()
	return User{}, nil
}

func (m *Model) Put(email, password string) error {
	return nil
}

func hashAndSalt(password string) (string, error) {
	// TODO: Use something other than MinCost?
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hash/salt password: %v", err)
	}

	return string(hash), nil
}
