package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

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

/*
This Authentication middleware usage is for route

mux.Handle("/v1/product/", middleware.Authentication(http.StripPrefix("/v1/product", app.ProductRouter())))
*/
func Authentication() gin.HandlerFunc {
	jwtSecret := os.Getenv("SECRET_KEY")

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Missing Authorization Header",
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid Token",
			})
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid Token Claims",
			})
			return
		}

		// Extract data from token
		c.Set(UserContextKey, UserClaims{
			UserID:  uint(claims["user_id"].(float64)),
			Email:   claims["email"].(string),
			IsAdmin: claims["is_admin"].(bool),
		})

		c.Next()
	}
}

/*
This AdminHandler usage is for each endpoint

authRouter.HandleFunc("POST /login", middleware.AdminHandler(app.login))
*/
func AdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get(UserContextKey)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": PleaseReLogin,
			})
			return
		}

		user := claims.(UserClaims)

		if !user.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": NotAuthorized,
			})
			return
		}

		c.Next()
	}
}

/*
This UserHandler usage is for each endpoint

authRouter.HandleFunc("POST /login", middleware.UserHandler(app.login))
*/
func UserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get(UserContextKey)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": PleaseReLogin,
			})
			return
		}

		user := claims.(UserClaims)

		if user.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": AdminNotAllowed,
			})
			return
		}

		c.Next()
	}
}
