package device

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

type Model struct {
	env.Env
}

func (m *Model) Create(u *env.User, name string, ctx context.Context) (*env.Device, error) {
	var err error

	// Creating device

	k := datastore.IncompleteKey("Device", nil)
	e := &env.Device{
		User: u.Key,
		Name: name,
	}

	e.Key, err = m.GetClient().Put(ctx, k, e)
	if err != nil {
		return nil, fmt.Errorf("Couldn't put new device in Datastore: %v", err)
	}

	// Linking device to user

	if u.Devices == nil {
		u.Devices = make([]*datastore.Key, 0, 1)
	}
	u.Devices = append(u.Devices, e.Key)

	_, err = m.GetClient().Mutate(ctx, datastore.NewUpdate(u.Key, u))
	if err != nil {
		// TODO: delete created device from DB?
		return nil, fmt.Errorf("Couldn't link device to user: %v", err)
	}

	return e, nil
}

func (m *Model) CreateTemp(ctx context.Context, k *datastore.Key, t env.Temp) error {
	var d env.Device

	err := m.GetClient().Get(ctx, k, &d)
	if err != nil {
		return fmt.Errorf("Key didn't match an existing entity: %v", err)
	}

	if d.History == nil {
		d.History = make([]env.Temp, 0, 1)
	}
	// TODO: this is pretty inefficient. Find another way of storing this data?
	// Push front
	d.History = append([]env.Temp{t}, d.History...)

	// Remove old temperatures
	if d.History[len(d.History)-1].Date.Unix() < time.Now().AddDate(0, 0, -14).Unix() {
		d.History = d.History[:len(d.History)-1]
	}

	_, err = m.GetClient().Mutate(ctx, datastore.NewUpdate(d.Key, &d))
	if err != nil {
		return fmt.Errorf("Couldn't store new temperature: %v", err)
	}

	return nil
}

func (m *Model) Alarm(ctx context.Context, key *datastore.Key, alarm *int) error {
	var d env.Device

	err := m.GetClient().Get(ctx, key, &d)
	if err != nil {
		return fmt.Errorf("Key didn't match an existing entity: %v", err)
	}

	d.Alarm = alarm

	_, err = m.GetClient().Mutate(ctx, datastore.NewUpdate(d.Key, &d))
	if err != nil {
		return fmt.Errorf("Couldn't update alarm: %v", err)
	}

	return nil
}
