package middleware

import (
	"bcc-project-v/sdk/config"
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
		sellerClaims, _ := c.Get("seller")
		sellerFromToken := sellerClaims.(model.SellerClaims)
		sellerCheck := entities.Seller{}
		if err := r.db.Debug().First(&sellerCheck, sellerFromToken.ID).Error; err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, "Can't find seller. Please try again!", nil)
			return
		}
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

		c.Set("seller", seller)
		c.Next()
	}
}

func (r *Repository) IsUserLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, _ := c.Get("user")
		userFromToken := userClaims.(model.UserClaims)
		userCheck := entities.User{}
		if err := r.db.Debug().First(&userCheck, userFromToken.ID).Error; err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, "Can't find user. Please try again!", nil)
			return
		}
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

		c.Set("user", user)
		c.Next()
	}
}
