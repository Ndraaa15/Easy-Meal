package handlers

import (
	"bcc-project-v/sdk/config"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/middleware"
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
	user := h.http.Group(v1.BasePath() + "/user")
	user.POST("/signup", h.UserRegister)
	user.GET("/login", h.UserLogin)
	user.Use(middleware.IsUserLoggedIn()).PUT("/market/update", h.UserUpdate)

	//Admin
	admin := h.http.Group(v1.BasePath() + "/admin")
	admin.POST("/signup", h.AdminRegister)
	admin.GET("/login", h.AdminLogin)
	admin.Use(middleware.IsAdminLoggedIn()).PUT("/update", h.AdminUpdate)

	//Product for admin
	product_admin := h.http.Group(v1.BasePath() + "/admin/market")
	product_admin.Use(middleware.IsAdminLoggedIn()).
		POST("/product/upload", h.PostProduct).
		POST("/product/upload/image", h.PostImageProduct).
		PUT("/product/:product_id", h.UpdateProduct).
		GET("/product", h.GetAllProduct).
		GET("/product/:product_id", h.GetProductByID)

	//Product for user
	product_user := h.http.Group(v1.BasePath() + "/user/market")
	product_user.Use(middleware.IsUserLoggedIn()).
		GET("/market/product", h.GetAllProduct).
		GET("/market/product/:product_id", h.GetProductByID)

}

func (h *handler) Run() {
	h.http.Run(fmt.Sprintf(":%s", h.config.GetEnv("PORT")))
}
