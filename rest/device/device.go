package device

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/auth"
	"github.com/NilsG-S/antifreeze-back-end/common/env"
	"github.com/NilsG-S/antifreeze-back-end/common/user"
)

func Apply(route *gin.RouterGroup, xEnv env.Env) {
	route.POST("/create", Create(xEnv))
	route.POST("/temp", auth.DeviceMiddleware(xEnv), Temp(xEnv))
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
	Date int64 `json:"date" binding:"required"`
	Temp int   `json:"temp" binding:"required"`
}

func Temp(xEnv env.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err    error
			json   TempInput
			dModel env.DeviceModel = xEnv.GetDevice()
		)

		// Binding data

		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invald input: %v", err),
			})
			return
		}

		// Decoding JWT

		dClaims := auth.GetDevice(c)

		var dKey *datastore.Key
		dKey, err = datastore.DecodeKey(dClaims.DeviceKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid device key: %v", err),
			})
			return
		}

		// Saving temp

		newT := env.Temp{
			Date:  json.Date,
			Value: json.Temp,
		}

		err = dModel.CreateTemp(c, dKey, &newT)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Unable to save temp in DB: %v", err),
			})
			return
		}

		xEnv.GetWS().PushTemp(dClaims.UserKey, dClaims.DeviceKey, &newT)

		c.Status(http.StatusOK)
	}
}
