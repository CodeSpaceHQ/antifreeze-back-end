package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

const (
	InvalidUsernamePassword = "Invalid username or password"
)

// TODO: split out token claim extraction for usage in ws server (getting email when authing)
// Regular function for usage in websocket server
func Verify(tokenString string, env *env.Env) error {
	// TODO: have this check scope of JWT (device vs user)
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

// Middleware for usage in Gin
func VerifyMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		// MUST PUT TOKEN IN `Authorization` HEADER
		token := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
		err := Verify(token, env)
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
