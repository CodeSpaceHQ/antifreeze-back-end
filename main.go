// contains the main "command" (running) logic
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"

	aCommon "github.com/NilsG-S/antifreeze-back-end/common/auth"
	dCommon "github.com/NilsG-S/antifreeze-back-end/common/device"
	uCommon "github.com/NilsG-S/antifreeze-back-end/common/user"
	aRoutes "github.com/NilsG-S/antifreeze-back-end/rest/auth"
	dRoutes "github.com/NilsG-S/antifreeze-back-end/rest/device"
	uRoutes "github.com/NilsG-S/antifreeze-back-end/rest/user"
	"github.com/NilsG-S/antifreeze-back-end/ws"
)

func main() {
	var err error

	// Setting up server "environment"

	env := &Env{
		Secret: os.Getenv("ANTIFREEZE_SECRET"),
	}

	env.Auth = &aCommon.Model{Env: env}
	env.Device = &dCommon.Model{Env: env}
	env.User = &uCommon.Model{Env: env}

	// Setting up logger

	var out *os.File = os.Stdout
	if e := os.Getenv("ANTIFREEZE_ENV"); e == "prod" {
		// Getting the current user to find their home directory
		var usr *user.User
		usr, err = user.Current()
		if err != nil {
			fmt.Printf("Couldn't get current user: %v", err)
			return
		}

		// Opening a log file in the current user's home directory
		out, err = os.OpenFile(usr.HomeDir+"/out.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("Couldn't open logfile: %v", err)
			return
		}
		defer out.Close()
	}
	env.Logger = log.New(out, "", log.LstdFlags|log.Lshortfile)

	// Setting up datastore client

	ctx := context.Background()
	// $DATASTORE_PROJECT_ID is used when second arg is empty
	// $GOOGLE_APPLICATION_CREDENTIALS points to credentials JSON
	env.Client, err = datastore.NewClient(ctx, "")
	if err != nil {
		fmt.Printf("Couldn't create client: %v", err)
		return
	}

	// Setting up router

	router := gin.New()

	router.Use(gin.LoggerWithWriter(out))
	router.Use(gin.Recovery())

	// Setting up server

	httpServer := &http.Server{
		Addr:           ":8081",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	env.WS = ws.NewServer(env)

	// Setting up routes

	routes(router, env)

	// Running Server

	go env.WS.RunServer()
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("ListenAndServe error: %v", err)
	}

	return
}

func routes(router *gin.Engine, env *Env) {
	// TODO: Add a NoRoute handler

	// TODO: Remove this once a real front-end exists
	router.StaticFile("/", "home.html")

	// # RESTful routes

	rest := router.Group("/rest")

	// ## User routes

	user := rest.Group("/user")
	uRoutes.Apply(user, env)

	// ## Auth routes

	auth := rest.Group("/auth")
	aRoutes.Apply(auth, env)

	// ## Device routes

	device := rest.Group("/device")
	dRoutes.Apply(device, env)

	// * WebSocket routes

	router.GET("/ws", func(c *gin.Context) {
		err := env.GetWS().Register(c.Writer, c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Couldn't register connection: %v", err),
			})
			return
		}

		c.Status(http.StatusOK)
	})
}
