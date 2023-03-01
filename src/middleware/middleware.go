package middleware

import (
	"bcc-project-v/sdk/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func IsAdminLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil header Authorization dari request
		authHeader := c.GetHeader("Authorization")
		// split string header menjadi 2 bagian (Bearer dan token)
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// ambil token dari bagian kedua header
		tokenString := tokenParts[1]

		// parse dan verifikasi token menggunakan secret key yang sama dengan saat membuat token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secretKey := config.Init().GetEnv("SECRET_KEY") // ambil secret key dari config
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// token valid, lanjutkan ke handler selanjutnya
		c.Next()
	}
}
