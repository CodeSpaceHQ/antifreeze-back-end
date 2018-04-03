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
	UserType                = "user"
	DeviceType              = "device"
)

type UserClaims struct {
	Type    string `json:"type"`
	UserKey string `json:"user_key"`
	jwt.StandardClaims
}

func (u *UserClaims) Valid() error { return nil }

type DeviceClaims struct {
	Type      string `json:"type"`
	UserKey   string `json:"user_key"`
	DeviceKey string `json:"device_key"`
}

func (d *DeviceClaims) Valid() error { return nil }

type Interface interface {
	Generate(jwt.Claims) (string, error)
	GetSecret() string
}

type Model struct {
	*env.Env
}

func (m *Model) Generate(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.GetSecret()))
}

// TODO: split out token claim extraction for usage in ws server (getting email when authing)
// Regular function for usage in websocket server
func Verify(tokenString string, env *env.Env) error {
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

	var tType string
	tType, ok = claims["type"].(string)
	if !ok {
		return fmt.Errorf("Token `type` wasn't a string")
	}

	if tType == "user" {
		ok = claims.VerifyExpiresAt(time.Now().Unix(), true)
		if !ok {
			return fmt.Errorf("Token has expired")
		}
	}

	return nil
}

func (m *Model) Decode(tString string, claims jwt.Claims) (*jwt.Token, error) {
	// TODO: doesn't the claims object passed here just get populated automatically?
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.GetSecret()), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("Invalid token: %v", err)
	}

	return token, nil
}

func (m *Model) DecodeUser(tString string) (*UserClaims, error) {
	token, err := Decode(tString, &UserClaims{})
	if err != nil {
		return nil, fmt.Errorf("Unable to decode UserClaims: %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("Claims weren't of type UserClaims")
	}

	return claims, nil
}

func (m *Model) DecodeDevice(tString string) (*DeviceClaims, error) {
	token, err := Decode(tString, &DeviceClaims{})
	if err != nil {
		return nil, fmt.Errorf("Unable to decode DeviceClaims: %v", err)
	}

	claims, ok := token.Claims.(*DeviceClaims)
	if !ok {
		return nil, fmt.Errorf("Claims weren't of type DeviceClaims")
	}

	return claims, nil
}

// TODO: UserMiddleware - checks for expiration date
// TODO: DeviceMiddleware

// TODO: Add scope checker for tokens? Long term, if ever

// Middleware for usage in Gin
func VerifyMiddleware(env *env.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Put decoded JWT into context?
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
