// contains the main "command" (running) logic
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common"
	"github.com/NilsG-S/antifreeze-back-end/common/db"
	"github.com/NilsG-S/antifreeze-back-end/ws"
)

func main() {
	var err error

	// Can be customized with gin.New()
	router := gin.Default()
	httpServer := &http.Server{
		Addr:           ":8081",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server := ws.NewServer()

	go server.RunServer()

	router.StaticFile("/", "home.html")

	router.POST("/test/post", func(c *gin.Context) {
		cur, err := db.GetInstance()
		if err != nil {
			c.String(http.StatusBadRequest, "Couldn't test posting: %v", err)
			return
		}

		cur.Testing()
		c.String(http.StatusOK, "Test posing successful!")
	})

	router.POST("/test/get", func(c *gin.Context) {
		cur, err := db.GetInstance()
		if err != nil {
			c.String(http.StatusBadRequest, "Couldn't test get: %v", err)
			return
		}

		cur.TestingGet()
		c.String(http.StatusOK, "Test get successful!")
	})

	router.Any("/ws", func(c *gin.Context) {
		server.Register(c.Writer, c.Request)
	})

	router.POST("/user/devices", func(c *gin.Context) {
		// TODO: This is a stopgap
		server.POSTUserDevices(1, "test@ttu.edu")
	})

	router.POST("/device/history", func(c *gin.Context) {
		mes := common.Temperature{
			Sub:      "/device/history",
			Op:       common.Add,
			DeviceID: 1,
			Temp:     32,
			Time:     time.Now(),
		}

		server.POSTDeviceHistory(mes)
	})

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}
