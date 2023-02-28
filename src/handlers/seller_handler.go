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
	newSeller := model.RegisterSeller{}
	if err := c.ShouldBindJSON(&newSeller); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	seller := entities.Seller{
		Shop:     newSeller.Shop,
		Email:    newSeller.Email,
		Password: newSeller.Password,
		Address:  newSeller.Address,
		Contact:  newSeller.Contact,
	}

	if err := h.Repository.CreateSeller(&seller); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't create new seller", nil)
	}
}

func (h *handler) SellerLogin(c *gin.Context) {
	sellerLogin := model.LoginSeller{}
	if err := c.ShouldBindJSON(&sellerLogin); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	sellerFound, err := h.Repository.FindSellerByEmail(sellerLogin)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the seller", nil)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(sellerFound.Password), []byte(sellerLogin.Password)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't compare the password", nil)
		return
	}

	//JWT TOKEN
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sellerFound.Shop,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	//GET JWT TOKEN
	tokenString, err := token.SignedString([]byte(sellerFound.Password))

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create token", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login Successful", nil)
	c.JSON(http.StatusOK, gin.H{
		"token JWT": tokenString,
	})

}

func (h *handler) SellerUpdate(c *gin.Context) {
	idReq := model.GetSellerByID{}
	updateSeller := model.UpdateSeller{}

	if err := c.ShouldBindUri(&idReq); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	if err := c.ShouldBindJSON(&updateSeller); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	sellerFound, err := h.Repository.FindSellerByID(idReq.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the seller", nil)
	}

	if updateSeller.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(updateSeller.Password), 20)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create password", nil)
			return
		}
		sellerFound.Password = string(hashPassword)
	} else if updateSeller.Shop != "" {
		sellerFound.Shop = updateSeller.Shop
	} else if updateSeller.Email != "" {
		sellerFound.Email = updateSeller.Email
	} else if updateSeller.Contact != "" {
		sellerFound.Contact = updateSeller.Contact
	} else if updateSeller.Address != "" {
		sellerFound.Address = updateSeller.Address
	}

	if err := h.Repository.UpdateSeller(sellerFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed updated seller", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &sellerFound)
}
