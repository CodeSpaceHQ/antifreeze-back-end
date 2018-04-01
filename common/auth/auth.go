package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

const (
	InvalidUsernamePassword = "Invalid username or password"
)

// Regular function for usage in websocket server
func Verify(tokenString string, env *env.Env) error {
	// TODO: have this check scope of JWT (device vs user)
	// time.Now().Unix() < jwt.exp
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(env.GetSecret()), nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("Invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("Couldn't map claims")
	}

	ok = claims.VerifyExpiresAt(time.Now().Unix(), true)
	if !ok {
		return fmt.Errorf("Token has expired")
	}

	return nil
}

type VerifyInput struct {
	Token string `json:"token" binding:"required"`
}

// Middleware for usage in Gin
func VerifyMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err  error
			json VerifyInput
		)

		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid input: %v", err),
			})
			return
		}

		err = Verify(json.Token, env)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid token: %v", err),
			})
			return
		}

		c.Next()
	}
}

// TODO: make generic version of JWT endpoint for user and device, then specialize
