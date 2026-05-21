package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID  uint   `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

const (
	UserContextKey  string = "user"
	PleaseReLogin   string = "Unauthorized, please re login!"
	NotAuthorized   string = "You are not authorized to perform this action!"
	AdminNotAllowed string = "Admins are not allowed to perform this action!"
)

func GetUserFromGin(c *gin.Context) (UserClaims, bool) {
	raw, exists := c.Get(UserContextKey)
	if !exists {
		return UserClaims{}, false
	}
	user, ok := raw.(UserClaims)
	return user, ok
}

func Authentication() gin.HandlerFunc {
	jwtSecret := os.Getenv("SECRET_KEY")

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
				Status:  http.StatusUnauthorized,
				Message: PleaseReLogin,
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
				Status:  http.StatusUnauthorized,
				Message: "Invalid Token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.BaseResponse{
				Status:  http.StatusUnauthorized,
				Message: "Invalid Token Claims",
			})
			return
		}

		role, _ := claims["role"].(string)
		c.Set(UserContextKey, UserClaims{
			UserID:  uint(claims["user_id"].(float64)),
			Email:   claims["email"].(string),
			IsAdmin: role == "admin",
		})

		c.Next()
	}
}
