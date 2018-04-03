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

func (m *Model) CreateTemp(ctx context.Context, k *datastore.Key, t Temp) error {
	var d Device

	err := m.Get(ctx, k, &d)
	if err != nil {
		return fmt.Errorf("Key didn't match an existing entity: %v", err)
	}

	if d.History == nil {
		d.History = make([]Temp, 0, 1)
	}
	// TODO: this is pretty inefficient. Find another way of storing this data?
	// Push front
	d.History = append([]Temp{t}, d.History...)

	// Remove old temperatures
	if d.History[len(d.History)-1].Date.Unix() < time.Now().AddDate(0, 0, -14).Unix() {
		d.History = d.History[:len(d.History)-1]
	}

	_, err = m.Mutate(ctx, datastore.NewUpdate(d.Key, d))
	if err != nil {
		return fmt.Errorf("Couldn't store new temperature: %v", err)
	}

	return nil
}
