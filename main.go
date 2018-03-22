// contains the main "command" (running) logic
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common"
	"github.com/NilsG-S/antifreeze-back-end/common/db"
	"github.com/NilsG-S/antifreeze-back-end/ws"
)

func main() {
	var (
		err error
		cur *db.Conn
	)

	out, err := os.OpenFile("/out.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer out.Close()
	log.SetOutput(out)

	// Can be customized with gin.New()
	router := gin.New()

	router.Use(gin.LoggerWithWriter(out))

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
		log.Println("Entering /test/post")

		cur, err = db.GetInstance()
		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Couldn't get DB connection: %v", err))
			return
		}

		cur.Testing()
		c.String(http.StatusOK, "Test posing successful!")
	})

	router.GET("/test/get", func(c *gin.Context) {
		log.Println("Entering /test/get")

		cur, err = db.GetInstance()
		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Couldn't get DB connection: %v", err))
			return
		}

		var response []string
		response, err = cur.TestingGet()
		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Couldn't test get: %v", err))
			return
		}

		for _, v := range response {
			log.Println(v)
		}

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
