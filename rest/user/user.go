package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

func Apply(route *gin.RouterGroup, env env.Env) {
	route.POST("/create", Create(env))
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
