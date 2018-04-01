package device

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/auth"
	"github.com/NilsG-S/antifreeze-back-end/common/device"
	"github.com/NilsG-S/antifreeze-back-end/common/env"
	"github.com/NilsG-S/antifreeze-back-end/common/user"
)

func Apply(route *gin.RouterGroup, env *env.Env) {
	uModel := &user.Model{Env: env}
	dModel := &device.Model{Env: env}

	route.POST("/create", Create(uModel, dModel))
}

type CreateInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

func Create(uModel user.Interface, dModel device.Interface) func(c *gin.Context) {
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

		exp := time.Now().AddDate(1, 0, 0).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			// TODO: These tokens shouldn't ever expire
			"exp":        exp,
			"device_key": d.Key.Encode(),
			"user_key":   u.Key.Encode(),
		})

		var tokenStr string
		tokenStr, err = token.SignedString([]byte(uModel.GetSecret()))
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
