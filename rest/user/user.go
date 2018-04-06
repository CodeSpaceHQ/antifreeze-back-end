package user

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/auth"
	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

func Apply(route *gin.RouterGroup, env env.Env) {
	route.POST("/create", Create(env))
	route.GET("/devices", auth.UserMiddleware(env), Devices(env))
}

type CreateInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Create(xEnv env.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err   error
			json  CreateInput
			model env.UserModel = xEnv.GetUser()
		)

		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invald input: %v", err),
			})
			return
		}

		err = model.Create(json.Email, json.Password, c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Couldn't create user: %v", err),
			})
			return
		}

		c.Status(http.StatusOK)
	}
}

func Devices(xEnv env.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err    error
			uModel env.UserModel = xEnv.GetUser()
		)

		// Decoding JWT

		uClaims := auth.GetUser(c)

		var uKey *datastore.Key
		uKey, err = datastore.DecodeKey(uClaims.UserKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid user key: %v", err),
			})
			return
		}

		// Get user

		var u *env.User
		u, err = uModel.GetByKey(c, uKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Couldn't get user by key: %v", err),
			})
			return
		}

		// Get user devices

		var d []env.GetDevicesJSON
		d, err = uModel.GetDevices(c, u)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Couldn't get user devices: %v", err),
			})
			return
		}

		// Automatically converts the struct to JSON
		c.JSON(http.StatusOK, gin.H{
			"devices": d,
		})
	}
}
