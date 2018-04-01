package device

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
	"github.com/NilsG-S/antifreeze-back-end/common/user"
)

type Temp struct {
	Value int       `datastore:"value,noindex"`
	Date  time.Time `datastore:"date,noindex"`
}

type Device struct {
	Key     *datastore.Key `datastore:"__key__"`
	Name    string         `datastore:"name,noindex"`
	Alarm   int            `datastore:"alarm,noindex"`
	User    *datastore.Key `datastore:"user,noindex"`
	History []Temp         `datastore:"history,noindex"`
}

// TODO: when storing time, remove from the end of the list if it's greater than 2 weeks ago
// Use Unix time

type Interface interface {
	Create(*user.User, string, context.Context) (*Device, error)
}

type Model struct {
	*env.Env
}

func (m *Model) Create(u *user.User, name string, ctx context.Context) (*Device, error) {
	var err error

	// Creating device

	k := datastore.IncompleteKey("Device", nil)
	e := &Device{
		User: u.Key,
		Name: name,
	}

	e.Key, err = m.Put(ctx, k, e)
	if err != nil {
		return nil, fmt.Errorf("Couldn't put new device in Datastore: %v", err)
	}

	// Linking device to user

	if u.Devices == nil {
		u.Devices = make([]*datastore.Key, 0, 1)
	}
	u.Devices = append(u.Devices, e.Key)

	_, err = m.Mutate(ctx, datastore.NewUpdate(u.Key, u))
	if err != nil {
		// TODO: delete created device from DB?
		return nil, fmt.Errorf("Couldn't link device to user: %v", err)
	}

	return e, nil
}
