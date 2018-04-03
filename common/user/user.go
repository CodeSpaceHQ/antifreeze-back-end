package user

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

// Make sure to handle case where no devices are present
type User struct {
	Key      *datastore.Key   `datastore:"__key__"`
	Email    string           `datastore:"email"`
	Password string           `datastore:"password,noindex"`
	Devices  []*datastore.Key `datastore:"devices,noindex"`
}

// In case we want to mock the model for unit tests
type Interface interface {
	// context.Context is an interface, so it shouldn't be a pointer anyway
	GetByEmail(string, context.Context) (*User, error)
	Create(string, string, context.Context) error
}

type Model struct {
	*env.Env
}

func (m *Model) GetByEmail(email string, ctx context.Context) (*User, error) {
	var u User

	q := datastore.NewQuery("User").Filter("email =", email)
	t := m.Run(ctx, q)
	_, err := t.Next(&u)

	if err == iterator.Done {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Error when iterating GetByEmail query: %v", err)
	}

	return &u, nil
}

func (m *Model) Create(email, password string, ctx context.Context) error {
	user, err := m.GetByEmail(email, ctx)
	if err != nil {
		return fmt.Errorf("Error when checking whether email already exists: %v", err)
	}
	if user != nil {
		return fmt.Errorf("Email already exists")
	}

	hash, err := hashAndSalt(password)
	if err != nil {
		return fmt.Errorf("Unable to hash/salt password: %v", err)
	}

	k := datastore.IncompleteKey("User", nil)
	e := &User{
		Email:    email,
		Password: hash,
	}

	_, err = m.Put(ctx, k, e)
	if err != nil {
		return fmt.Errorf("Couldn't put new user in Datastore, %v", err)
	}

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

func ComparePassword(hashed, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return fmt.Errorf("Password comparison failed: %v", err)
	}

	return nil
}
