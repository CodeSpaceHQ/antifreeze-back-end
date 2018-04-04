package device

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/auth"
	"github.com/NilsG-S/antifreeze-back-end/common/env"
	"github.com/NilsG-S/antifreeze-back-end/common/user"
)

func Apply(route *gin.RouterGroup, env env.Env) {
	route.POST("/create", Create(env))
}

type CreateInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

func Create(xEnv env.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err    error
			json   CreateInput
			aModel env.AuthModel   = xEnv.GetAuth()
			dModel env.DeviceModel = xEnv.GetDevice()
			uModel env.UserModel   = xEnv.GetUser()
		)

		// Binding data

		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invald input: %v", err),
			})
			return
		}

		// Getting user

		var u *env.User
		u, err = uModel.GetByEmail(json.Email, c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Couldn't find user: %v", err),
			})
			return
		}
		if u == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": auth.InvalidUsernamePassword,
			})
			return
		}

		// Password comparison

		err = user.ComparePassword(u.Password, json.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": auth.InvalidUsernamePassword,
			})
			return
		}

		// Creating device

		var d *env.Device
		d, err = dModel.Create(u, json.Name, c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Couldn't create device: %v", err),
			})
			return
		}

		// Making JWT

		var tokenStr string
		tokenStr, err = aModel.Generate(&env.DeviceClaims{
			Type:      auth.DeviceType,
			UserKey:   u.Key.Encode(),
			DeviceKey: d.Key.Encode(),
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Unable to generate token: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenStr,
		})
	}
}

type TempInput struct {
	Time int64 `json:"time" binding:"required"`
	Temp int   `json:"temp" binding:"required"`
}

func Temp(xEnv env.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err  error
			json TempInput
			// aModel env.AuthModel   = xEnv.GetAuth()
			// dModel env.DeviceModel = xEnv.GetDevice()
		)

		// Binding data

		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invald input: %v", err),
			})
			return
		}

		// Decoding JSON

		// dClaims := auth.GetDevice(c)

		// TODO: Save temp to device
		// TODO: Push temp to user

		c.Status(http.StatusOK)
	}
}
