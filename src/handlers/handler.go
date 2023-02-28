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

	//User
	v1.POST("/user/signup", h.UserRegister)
	v1.GET("/user/login", h.UserLogin)
	v1.PUT("/user/update/:id", h.UserUpdate)

	//Admin
	v1.POST("/admin/signup", h.AdminRegister)
	v1.GET("/admin/login", h.AdminLogin)
	v1.PUT("/admin/update/:id", h.AdminUpdate)

}

func (h *handler) Run() {
	h.http.Run(fmt.Sprintf(":%s", h.config.GetEnv("PORT")))
}
