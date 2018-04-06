package user

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

type Model struct {
	env.Env
}

func (m *Model) GetByEmail(email string, ctx context.Context) (*env.User, error) {
	var u env.User

	q := datastore.NewQuery("User").Filter("email =", email)
	t := m.GetClient().Run(ctx, q)
	_, err := t.Next(&u)

	if err == iterator.Done {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Error when iterating GetByEmail query: %v", err)
	}

	return &u, nil
}

func (m *Model) GetByKey(ctx context.Context, key *datastore.Key) (*env.User, error) {
	var u env.User

	err := m.GetClient().Get(ctx, key, &u)
	if err != nil {
		return nil, fmt.Errorf("Key didn't match an existing user: %v", err)
	}

	return &u, nil
}

func (m *Model) GetDevices(ctx context.Context, user *env.User) ([]env.GetDevicesJSON, error) {
	devices := make([]env.Device, len(user.Devices))
	err := m.GetClient().GetMulti(ctx, user.Devices, devices)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get devices: %v", err)
	}

	out := make([]env.GetDevicesJSON, len(user.Devices))
	for i, v := range devices {
		out[i] = env.GetDevicesJSON{
			DeviceKey: v.Key.Encode(),
			Name:      v.Name,
			Alarm:     v.Alarm,
		}
	}

	return out, nil
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
	e := &env.User{
		Email:    email,
		Password: hash,
	}

	_, err = m.GetClient().Put(ctx, k, e)
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
