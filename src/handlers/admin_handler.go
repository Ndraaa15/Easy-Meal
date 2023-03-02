package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *handler) AdminRegister(c *gin.Context) {
	newAdmin := model.AdminRegister{}
	if err := c.ShouldBindJSON(&newAdmin); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.Password), 12)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Password format is incorrect. Please follow the specified format and try again!", nil)
		return
	}

	admin := entities.Admin{
		Shop:     newAdmin.Shop,
		Email:    newAdmin.Email,
		Password: string(hashPassword),
	}

	if err := h.Repository.CreateAdmin(&admin); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't create new admin", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Register successful! Welcome, "+admin.Shop+"!", admin)
}

func (h *handler) AdminLogin(c *gin.Context) {
	adminLogin := model.AdminLogin{}
	if err := c.ShouldBindJSON(&adminLogin); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	adminFound, err := h.Repository.FindAdminByEmail(&adminLogin)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Admin not found. Please try again with a valid email!", nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(adminFound.Password), []byte(adminLogin.Password)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Wrong password. Please try again with a valid password!", nil)
		return
	}

	//JWT TOKEN
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   adminFound.Shop,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"shop":  adminFound.Shop,
		"email": adminFound.Email,
		"id":    adminFound.ID,
	})

	//GET JWT TOKEN
	tokenString, err := token.SignedString([]byte(h.config.GetEnv("SECRET_KEY")))

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create token JWT. Please try again to login!", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login successful! Welcome back, "+adminFound.Shop+"!", nil)
	c.JSON(http.StatusOK, gin.H{
		"JWT Token": tokenString,
	})

}

func (h *handler) AdminUpdate(c *gin.Context) {
	adminClaims, _ := c.Get("admin")
	admin := adminClaims.(model.AdminClaims)
	updateAdmin := model.AdminUpdate{}

	if err := c.ShouldBindJSON(&updateAdmin); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	adminFound, err := h.Repository.FindAdminByID(admin.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Admin not found. Please try again later!", nil)
		return
	}

	if updateAdmin.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(updateAdmin.Password), 12)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create new password", nil)
			return
		}
		adminFound.Password = string(hashPassword)
	}
	if updateAdmin.Shop != "" {
		adminFound.Shop = updateAdmin.Shop
	}
	if updateAdmin.Email != nil {
		adminFound.Email = updateAdmin.Email
	}

	if err := h.Repository.UpdateAdmin(adminFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &adminFound)
}
