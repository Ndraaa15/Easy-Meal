package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"net/http"
	"time"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *handler) SellerRegister(c *gin.Context) {
	newSeller := model.SellerRegister{}
	if err := c.ShouldBindJSON(&newSeller); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", err.Error())
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newSeller.Password), 12)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Password format is incorrect. Please follow the specified format and try again!", err.Error())
		return
	}

	seller := entities.Seller{
		Shop:     newSeller.Shop,
		Email:    newSeller.Email,
		Password: string(hashPassword),
		Address:  newSeller.Address,
		Contact:  newSeller.Contact,
	}
	if err := h.Repository.CreateSeller(&seller); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed add seller to database. Please try again later or contact customer service for help", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Register successful! Welcome, "+seller.Shop+"!", seller)
}

func (h *handler) SellerLogin(c *gin.Context) {
	sellerLogin := model.SellerLogin{}
	if err := c.ShouldBindJSON(&sellerLogin); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", err.Error())
		return
	}

	sellerFound, err := h.Repository.FindSellerByEmail(&sellerLogin)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Seller not found. Please try again with a valid email!", err.Error())
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(sellerFound.Password), []byte(sellerLogin.Password)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Wrong password. Please try again with a valid password!", err.Error())
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   sellerFound.Shop,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"id":    sellerFound.ID,
		"email": sellerFound.Email,
	})

	tokenString, err := token.SignedString([]byte(h.config.GetEnv("SECRET_KEY")))

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create token JWT. Please try again to login!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login successful! Welcome back, "+sellerFound.Shop+" ! ", tokenString)
}

func (h *handler) SellerUpdate(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)

	shop := c.PostForm("shop")
	email := c.PostForm("email")
	address := c.PostForm("address")
	contact := c.PostForm("contact")

	supClient := supabasestorageuploader.NewSupabaseClient(
		"https://arcudskzafkijqukfool.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFyY3Vkc2t6YWZraWpxdWtmb29sIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Nzc2NDk3MjksImV4cCI6MTk5MzIyNTcyOX0.CjOVpoFAdq3U-AeAzsuyV6IGcqx2ZnaXjneTis5qd6w",
		"bcc-project",
		"seller-image",
	)
	file, err := c.FormFile("seller_image")
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to get seller image", err.Error())
		return
	}
	link, err := supClient.Upload(file)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload seller image", err.Error())
		return
	}

	sellerFound, err := h.Repository.FindSellerByID(seller.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Seller not found. Please try again later!", err.Error())
		return
	}

	if shop != "" {
		sellerFound.Shop = shop
	}
	if email != "" {
		sellerFound.Email = email
	}
	if address != "" {
		sellerFound.Address = address
	}
	if contact != "" {
		sellerFound.Contact = contact
	}
	if sellerFound.SellerImage != link {
		sellerFound.SellerImage = link
	}

	if err := h.Repository.UpdateSeller(sellerFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &sellerFound)
}

func (h *handler) SellerUpdatePassword(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)

	sellerFound, err := h.Repository.FindSellerByID(seller.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Seller not found. Please try again later!", err.Error())
		return
	}

	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	if err = bcrypt.CompareHashAndPassword([]byte(sellerFound.Password), []byte(oldPassword)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Wrong password. Please try again with a valid password!", err.Error())
		return
	}

	if newPassword != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create new password", err.Error())
			return
		}
		sellerFound.Password = string(hashPassword)
	}

	if err := h.Repository.UpdateSeller(sellerFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Update Password Successful", &sellerFound)
}

// func (h *handler) GetOrder(c *gin.Context){
// 	cart :=
// }
