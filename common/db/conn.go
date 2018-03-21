package db

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/datastore"
)

type Conn struct {
	context context.Context
	client  *datastore.Client
}

var instance *Conn
var once sync.Once

// TODO(NilsG-S): Don't use singletons
func GetInstance() (*Conn, error) {
	var err error

	once.Do(func() {
		var (
			ctx context.Context = context.Background()
			cli *datastore.Client
		)

		// TODO(NilsG-S): determine project name dev/prod?
		cli, err = datastore.NewClient(ctx, "antifreeze")
		if err != nil {
			err = fmt.Errorf("Couldn't create client: %v", err)
		}

		instance = &Conn{
			context: ctx,
			client:  cli,
		}
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

type Test struct {
	Value string
}

func (c *Conn) Testing() error {
	k := datastore.IncompleteKey("Test", nil)
	e := new(Test)
	e.Value = "Testing"

	if _, err := c.client.Put(c.context, k, e); err != nil {
		return fmt.Errorf("Couldn't insert test entity: %v", err)
	}

	return nil
}
