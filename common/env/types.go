package env

import (
	"cloud.google.com/go/datastore"
	"github.com/dgrijalva/jwt-go"
)

// Auth types

type UserClaims struct {
	Type    string `json:"type"`
	UserKey string `json:"user_key"`
	jwt.StandardClaims
}

func (u *UserClaims) Valid() error { return nil }

type DeviceClaims struct {
	Type      string `json:"type"`
	UserKey   string `json:"user_key"`
	DeviceKey string `json:"device_key"`
}

func (d *DeviceClaims) Valid() error { return nil }

// Device types

type Temp struct {
	Value int   `json:"value"`
	Date  int64 `json:"date"`
}

type Device struct {
	Key     *datastore.Key `datastore:"__key__"`
	Name    string         `datastore:"name,noindex"`
	Alarm   int            `datastore:"alarm,noindex"`
	User    *datastore.Key `datastore:"user,noindex"`
	History []string       `datastore:"history,noindex"`
}

// User types

type User struct {
	Key      *datastore.Key   `datastore:"__key__"`
	Email    string           `datastore:"email"`
	Password string           `datastore:"password,noindex"`
	Devices  []*datastore.Key `datastore:"devices,noindex"`
}
