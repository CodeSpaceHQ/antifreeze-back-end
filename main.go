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
	"github.com/NilsG-S/antifreeze-back-end/ws"
)

type Env struct {
	*datastore.Client
	*log.Logger
	*ws.Server
}

func routes(router *gin.Engine, env *Env) {
	// TODO: Add a NoRoute handler
	router.StaticFile("/", "home.html")

	router.Any("/ws", func(c *gin.Context) {
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

func main() {
	var err error

	// Setting up logger

	var usr *user.User
	usr, err = user.Current()
	if err != nil {
		fmt.Printf("Couldn't get current user: %v", err)
		return
	}

	// TODO: Have the output split between file and stdout
	// Or just have a production ENV variable that sets file/stdout
	out, err := os.OpenFile(usr.HomeDir+"/out.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Couldn't open logfile: %v", err)
		return
	}
	defer out.Close()
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

	env := &Env{
		cli,
		logger,
		server,
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
