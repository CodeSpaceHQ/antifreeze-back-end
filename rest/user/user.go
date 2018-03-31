package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
	"github.com/NilsG-S/antifreeze-back-end/common/user"
)

func Apply(route *gin.RouterGroup, env *env.Env) {
	model := &user.Model{Env: env}

	route.POST("/create", Create(model))
}

type CreateInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Confirm  string `json:"confirm" binding:"required"`
}

func Create(model user.Interface) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err  error
			json CreateInput
		)

		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invald input: %v", err),
			})
			return
		}

		if json.Password != json.Confirm {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Passwords didn't match",
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
