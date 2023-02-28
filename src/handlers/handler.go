package handlers

import (
	"bcc-project-v/sdk/config"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	http       *gin.Engine
	config     config.Initializer
	db         *gorm.DB
	Repository repository.Repository
}

func Init(config config.Initializer, repo *repository.Repository) *handler {
	rest := handler{
		http:       gin.Default(),
		config:     config,
		Repository: *repo,
	}
	rest.registerRoutes()
	return &rest
}

func (h *handler) registerRoutes() {
	repository.NewRepository(h.db)
	h.http.GET("/", func(ctx *gin.Context) {
		helper.SuccessResponse(ctx, http.StatusOK, "Hello World", nil)
	})

	v1 := h.http.Group("/api/v1")

	v1.POST("/user/signup", h.UserRegister)
	v1.GET("/user/login", h.UserLogin)
	v1.PUT("/user/update/:id", h.UserUpdate)

	v1.POST("/seller/signup", h.SellerRegister)
	v1.GET("/seller/login", h.SellerLogin)
	v1.PUT("/seller/update/:id", h.SellerUpdate)

}

func (h *handler) Run() {
	h.http.Run(fmt.Sprintf(":%s", h.config.GetEnv("PORT")))
}
