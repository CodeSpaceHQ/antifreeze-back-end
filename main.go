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

	"github.com/NilsG-S/antifreeze-back-end/common"
	"github.com/NilsG-S/antifreeze-back-end/common/env"
	userRoutes "github.com/NilsG-S/antifreeze-back-end/rest/user"
	"github.com/NilsG-S/antifreeze-back-end/ws"
)

func main() {
	var err error

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
	logger := log.New(out, "", log.LstdFlags|log.Lshortfile)

	// Setting up datastore client

	var cli *datastore.Client
	ctx := context.Background()
	// $DATASTORE_PROJECT_ID is used when second arg is empty
	// $GOOGLE_APPLICATION_CREDENTIALS points to credentials JSON
	cli, err = datastore.NewClient(ctx, "")
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
	server := ws.NewServer()

	// Setting up server "environment"

	env := &env.Env{
		Client: cli,
		Logger: logger,
		Server: server,
	}

	// Setting up routes

	routes(router, env)

	// Running Server

	go server.RunServer()
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("ListenAndServe error: %v", err)
	}

	return
}

func routes(router *gin.Engine, env *env.Env) {
	// TODO: Add a NoRoute handler

	// # RESTful routes

	rest := router.Group("/rest")

	// ## User routes

	user := rest.Group("/user")
	userRoutes.Apply(user, env)

	// Old routes

	router.StaticFile("/", "home.html")

	router.GET("/ws", func(c *gin.Context) {
		env.Register(c.Writer, c.Request)
	})

	router.POST("/user/devices", func(c *gin.Context) {
		// TODO: This is a stopgap
		env.POSTUserDevices(1, "test@ttu.edu")
	})

	router.POST("/device/history", func(c *gin.Context) {
		mes := common.Temperature{
			Sub:      "/device/history",
			Op:       common.Add,
			DeviceID: 1,
			Temp:     32,
			Time:     time.Now(),
		}

		env.POSTDeviceHistory(mes)
	})
}
