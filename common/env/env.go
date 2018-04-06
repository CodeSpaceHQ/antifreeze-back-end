// NOTE: This file should never import anything from inside this project
package env

import (
	"context"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/dgrijalva/jwt-go"
)

// Methods for auth model
type AuthModel interface {
	Generate(jwt.Claims) (string, error)
	Decode(string, jwt.Claims) (*jwt.Token, error)
	DecodeUser(string) (*UserClaims, error)
	DecodeDevice(string) (*DeviceClaims, error)
}

// Methods for device model
type DeviceModel interface {
	Create(*User, string, context.Context) (*Device, error)
	CreateTemp(ctx context.Context, key *datastore.Key, temp Temp) error
}

// Methods for user model
type UserModel interface {
	// context.Context is an interface, so it shouldn't be a pointer anyway
	GetByEmail(string, context.Context) (*User, error)
	GetByKey(ctx context.Context, key *datastore.Key) (*User, error)
	GetDevices(ctx context.Context, user *User) ([]GetDevicesJSON, error)
	Create(string, string, context.Context) error
}

// Methods for WS server
type WS interface {
	RunServer()
	Register(w http.ResponseWriter, r *http.Request)
	PushTemp(userKey, deviceKey string, temp Temp)
}

type Env interface {
	GetClient() *datastore.Client

	GetAuth() AuthModel
	GetDevice() DeviceModel
	GetUser() UserModel
	GetWS() WS

	GetSecret() string
}
