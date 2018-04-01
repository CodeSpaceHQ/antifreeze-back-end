package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/auth"
	"github.com/NilsG-S/antifreeze-back-end/common/env"
	"github.com/NilsG-S/antifreeze-back-end/common/user"
)

func Apply(route *gin.RouterGroup, env *env.Env) {
	model := &user.Model{Env: env}

	route.POST("/login", Login(model))
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(model user.Interface) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err  error
			json LoginInput
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
		u, err = model.GetByEmail(json.Email, c)
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

		// Making JWT

		// TODO: Consider adding type claim for device vs user
		exp := time.Now().AddDate(1, 0, 0).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": exp,
			// TODO: use User key encoding instead
			"email": json.Email,
		})

		var tokenStr string
		tokenStr, err = token.SignedString([]byte(model.GetSecret()))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Unable to generate token: %v", err),
			})
			return
		}

		// TODO: Add fully populated user with devices
		c.JSON(http.StatusOK, gin.H{
			"token": tokenStr,
		})
	}
}
