package env

import (
	"log"

	"cloud.google.com/go/datastore"
	"github.com/NilsG-S/antifreeze-back-end/ws"
)

type Env struct {
	*datastore.Client
	*log.Logger
	*ws.Server
}
