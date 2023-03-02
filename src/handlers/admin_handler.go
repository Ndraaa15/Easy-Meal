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
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.Password), 12)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create password", nil)
		return
	}

	admin := entities.Admin{
		Shop:     newAdmin.Shop,
		Email:    newAdmin.Email,
		Password: string(hashPassword),
	}

	if err := h.Repository.CreateAdmin(&admin); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't create new admin", nil)
	}

	helper.SuccessResponse(c, http.StatusOK, "Register Successful", admin)
}

func (h *handler) AdminLogin(c *gin.Context) {
	adminLogin := model.AdminLogin{}
	if err := c.ShouldBindJSON(&adminLogin); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	adminFound, err := h.Repository.FindAdminByEmail(&adminLogin)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the admin", nil)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(adminFound.Password), []byte(adminLogin.Password)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't compare the password", nil)
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
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create token", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login Successful", nil)
	c.JSON(http.StatusOK, gin.H{
		"jwtToken": tokenString,
	})

}

func (h *handler) AdminUpdate(c *gin.Context) {
	idReq := model.GetAdminByID{}
	updateAdmin := model.AdminUpdate{}

	if err := c.ShouldBindUri(&idReq); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	if err := c.ShouldBindJSON(&updateAdmin); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	adminFound, err := h.Repository.FindAdminByID(idReq.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the admin", nil)
	}

	if updateAdmin.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(updateAdmin.Password), 12)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create password", nil)
			return
		}
		adminFound.Password = string(hashPassword)
	} else if updateAdmin.Shop != "" {
		adminFound.Shop = updateAdmin.Shop
	} else if updateAdmin.Email != nil {
		adminFound.Email = updateAdmin.Email
	}

	if err := h.Repository.UpdateAdmin(adminFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed updated admin", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &adminFound)
}
