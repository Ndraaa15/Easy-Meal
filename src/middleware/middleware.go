package middleware

import (
	"bcc-project-v/sdk/config"
	"bcc-project-v/src/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) IsSellerLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		tokenString := tokenParts[1]
		seller := model.SellerClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &seller, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secretKey := config.Init().GetEnv("SECRET_KEY")
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		// Checking seller (?)
		// if err := r.db.Debug().First(&entities.Seller{}, seller.ID).Error; err != nil {
		// 	helper.ErrorResponse(c, http.StatusBadRequest, "Can't find seller. Please try again!", nil)
		// 	return
		// }
		c.Set("seller", seller)
		c.Next()
	}
}

func (r *Repository) IsUserLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenParts := strings.SplitN(authHeader, " ", 2)

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		tokenString := tokenParts[1]
		user := model.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &user, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secretKey := config.Init().GetEnv("SECRET_KEY")
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		// Checking user (?)
		// if err := r.db.Debug().First(&entities.User{}, &user.ID).Error; err != nil {
		// 	helper.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		// 	return
		// }

		c.Set("user", user)
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowCredentials: true,
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
			// Enable Debugging for testing, consider disabling in production
			Debug: true,
		})
		c.HandlerFunc(ctx.Writer, ctx.Request)
		ctx.Next()
	}
}
