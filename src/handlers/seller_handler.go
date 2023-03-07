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

func (h *handler) SellerRegister(c *gin.Context) {
	newSeller := model.SellerRegister{}
	if err := c.ShouldBindJSON(&newSeller); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newSeller.Password), 12)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Password format is incorrect. Please follow the specified format and try again!", nil)
		return
	}

	seller := entities.Seller{
		Shop:     newSeller.Shop,
		Email:    newSeller.Email,
		Password: string(hashPassword),
	}

	if err := h.Repository.CreateSeller(&seller); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't create new seller", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Register successful! Welcome, "+seller.Shop+"!", seller)
}

func (h *handler) SellerLogin(c *gin.Context) {
	sellerLogin := model.SellerLogin{}
	if err := c.ShouldBindJSON(&sellerLogin); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	sellerFound, err := h.Repository.FindSellerByEmail(&sellerLogin)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Admin not found. Please try again with a valid email!", nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(sellerFound.Password), []byte(sellerLogin.Password)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Wrong password. Please try again with a valid password!", nil)
		return
	}

	//JWT TOKEN
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   sellerFound.Shop,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"shop":  sellerFound.Shop,
		"email": sellerFound.Email,
		"id":    sellerFound.ID,
	})

	//GET JWT TOKEN
	tokenString, err := token.SignedString([]byte(h.config.GetEnv("SECRET_KEY")))

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create token JWT. Please try again to login!", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login successful! Welcome back, "+sellerFound.Shop+"!", nil)
	c.JSON(http.StatusOK, gin.H{
		"JWT Token": tokenString,
	})

}

func (h *handler) SellerUpdate(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)
	updateAdmin := model.SellerUpdate{}

	if err := c.ShouldBindJSON(&updateAdmin); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}
	shop := c.PostForm("shop")
	email := c.PostForm("email")
	password := c.PostForm("password")
	address := c.PostForm("address")
	contact := c.PostForm("contact")
	sellerFound, err := h.Repository.FindSellerByID(seller.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Admin not found. Please try again later!", nil)
		return
	}

	if password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create new password", nil)
			return
		}
		sellerFound.Password = string(hashPassword)
	}
	if shop != "" {
		sellerFound.Shop = shop
	}
	if email != "" {
		sellerFound.Email = &email
	}
	if address != "" {
		sellerFound.Address = address
	}
	if contact != "" {
		sellerFound.Contact = contact
	}

	if err := h.Repository.UpdateSeller(sellerFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &sellerFound)
}
