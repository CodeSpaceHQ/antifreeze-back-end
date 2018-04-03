package device

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/auth"
	"github.com/NilsG-S/antifreeze-back-end/common/device"
	"github.com/NilsG-S/antifreeze-back-end/common/env"
	"github.com/NilsG-S/antifreeze-back-end/common/user"
)

func Apply(route *gin.RouterGroup, env *env.Env) {
	aModel := &auth.Model{Env: env}
	dModel := &device.Model{Env: env}
	uModel := &user.Model{Env: env}

	route.POST("/create", Create(uModel, dModel, aModel))
}

type CreateInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

func Create(uModel user.Interface, dModel device.Interface, aModel auth.Interface) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err  error
			json CreateInput
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

		var u *user.User
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

		var d *device.Device
		d, err = dModel.Create(u, json.Name, c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Couldn't create device: %v", err),
			})
			return
		}

		// Making JWT

		var tokenStr string
		tokenStr, err = aModel.Generate(&auth.DeviceClaims{
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
