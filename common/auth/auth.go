package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

const (
	InvalidUsernamePassword = "Invalid username or password"
	UserType                = "user"
	DeviceType              = "device"
	ClaimsKey               = "claims"
)

type Model struct {
	env.Env
}

func (m *Model) Generate(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.GetSecret()))
}

func (m *Model) Decode(tString string, claims jwt.Claims) (*jwt.Token, error) {
	// TODO: doesn't the claims object passed here just get populated automatically?
	// TODO: in that vein, should these functions return pointers or actual structs?
	token, err := jwt.ParseWithClaims(tString, claims, func(token *jwt.Token) (interface{}, error) {
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

// Pretty sure this already checks for token expiration
func (m *Model) DecodeUser(tString string) (*env.UserClaims, error) {
	token, err := m.Decode(tString, &env.UserClaims{})
	if err != nil {
		return nil, fmt.Errorf("Unable to decode UserClaims: %v", err)
	}

	claims, ok := token.Claims.(*env.UserClaims)
	if !ok {
		return nil, fmt.Errorf("Claims weren't of type UserClaims")
	}

	return claims, nil
}

func (m *Model) DecodeDevice(tString string) (*env.DeviceClaims, error) {
	token, err := m.Decode(tString, &env.DeviceClaims{})
	if err != nil {
		return nil, fmt.Errorf("Unable to decode DeviceClaims: %v", err)
	}

	claims, ok := token.Claims.(*env.DeviceClaims)
	if !ok {
		return nil, fmt.Errorf("Claims weren't of type DeviceClaims")
	}

	return claims, nil
}

// Middleware for usage in Gin
func UserMiddleware(env env.Env) gin.HandlerFunc {
	// TODO: Add scope checker for tokens? Long term, if ever
	model := env.GetAuth()

	return func(c *gin.Context) {
		// MUST PUT TOKEN IN `Authorization` HEADER
		token := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
		claims, err := model.DecodeUser(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid User token: %v", err),
			})
			return
		}

		c.Set(ClaimsKey, claims)

		c.Next()
	}
}

func DeviceMiddleware(env env.Env) gin.HandlerFunc {
	model := env.GetAuth()

	return func(c *gin.Context) {
		token := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
		claims, err := model.DecodeDevice(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid Device token: %v", err),
			})
			return
		}

		c.Set(ClaimsKey, claims)

		c.Next()
	}
}

func GetUser(c *gin.Context) *env.UserClaims {
	claims, exists := c.Get(ClaimsKey)
	if !exists {
		fmt.Println("Programmer Error (GetUser): claims not present")
		return nil
	}

	uClaims, ok := claims.(*env.UserClaims)
	if !ok {
		fmt.Println("Programmer Error (GetUser): claims should be *UserClaims")
		return nil
	}

	return uClaims
}

func GetDevice(c *gin.Context) *env.DeviceClaims {
	claims, exists := c.Get(ClaimsKey)
	if !exists {
		fmt.Println("Programmer Error (GetDevice): claims not present")
		return nil
	}

	dClaims, ok := claims.(*env.DeviceClaims)
	if !ok {
		fmt.Println("Programmer Error (GetDevice): claims should be *DeviceClaims")
		return nil
	}

	return dClaims
}
