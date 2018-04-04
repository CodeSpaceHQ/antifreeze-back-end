package main

import (
	"log"

	"cloud.google.com/go/datastore"

	"github.com/NilsG-S/antifreeze-back-end/common/auth"
	"github.com/NilsG-S/antifreeze-back-end/common/device"
	"github.com/NilsG-S/antifreeze-back-end/common/env"
	"github.com/NilsG-S/antifreeze-back-end/common/user"
	"github.com/NilsG-S/antifreeze-back-end/ws"
)

type Env struct {
	Client *datastore.Client
	Logger *log.Logger

	Auth   *auth.Model
	Device *device.Model
	User   *user.Model

	WS *ws.Server

	Secret string
}

func (e *Env) GetClient() *datastore.Client { return e.Client }

func (e *Env) GetAuth() env.AuthModel { return e.Auth }

func (e *Env) GetDevice() env.DeviceModel { return e.Device }

func (e *Env) GetUser() env.UserModel { return e.User }

func (e *Env) GetWS() env.WS { return e.WS }

func (e *Env) GetSecret() string { return e.Secret }
