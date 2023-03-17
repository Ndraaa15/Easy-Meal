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
		Username: newSeller.Username,
		Email:    newSeller.Email,
		Password: string(hashPassword),
		Address:  newSeller.Address,
		City:     newSeller.City,
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

	helper.SuccessResponse(c, http.StatusOK, tokenString, &sellerFound)
}

func (h *handler) SellerUpdate(c *gin.Context) {
	sellerClaims, exist := c.Get("seller")
	if !exist {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to load JWT token, please try again!", nil)
	}
	seller := sellerClaims.(model.SellerClaims)

	shop := c.PostForm("shop")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	address := c.PostForm("address")
	city := c.PostForm("city")
	contact := c.PostForm("contact")
	linkMaps := c.PostForm("link_maps")

	supClient := supabasestorageuploader.NewSupabaseClient(
		"https://arcudskzafkijqukfool.supabase.co",
		h.config.GetEnv("SUPABASE_API_KEY"),
		h.config.GetEnv("SUPABASE_STORAGE"),
		h.config.GetEnv("SUPABASE_SELLER_FOLDER"),
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
	if username != "" {
		sellerFound.Username = username
	}
	if email != "" {
		sellerFound.Email = email
	}
	if address != "" {
		sellerFound.Address = address
	}
	if city != "" {
		sellerFound.City = city
	}
	if contact != "" {
		sellerFound.Contact = contact
	}
	if linkMaps != "" {
		sellerFound.LinkMaps = linkMaps
	}
	if sellerFound.SellerImage != link {
		sellerFound.SellerImage = link
	}
	if password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create new password", err.Error())
			return
		}
		sellerFound.Password = string(hashPassword)
	}
	if err := h.Repository.UpdateSeller(sellerFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &sellerFound)
}

func (h *handler) GetOrder(c *gin.Context) {
	sellerClaims, exist := c.Get("seller")
	if !exist {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to load JWT token, please try again!", nil)
	}
	seller := sellerClaims.(model.SellerClaims)

	productsOrder, err := h.Repository.GetOrder(seller.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to get order product", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product Order Found!", &productsOrder)
}

func (h *handler) SetStatusSuccess(c *gin.Context) {
	checkOrder := model.CheckOrder{}
	if err := c.ShouldBindJSON(&checkOrder); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid data format", err.Error())
		return
	}

	payment, err := h.Repository.CheckOrder(checkOrder.PaymentCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to get payment", err.Error())
	}

	payment.StatusID = 2
	status := entities.Status{}
	if err := h.Repository.FindStatus(&status, 2); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed get status", err.Error())
		return
	}
	payment.Status = status

	if err := h.Repository.SavePayment(payment); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save udpated payment", err.Error())
	}
	helper.SuccessResponse(c, http.StatusOK, "Update payment Successful", &payment)
}

func (h *handler) SetStatusFailed(c *gin.Context) {
	checkOrder := model.CheckOrder{}
	if err := c.ShouldBindJSON(&checkOrder); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid data format", err.Error())
		return
	}

	payment, err := h.Repository.CheckOrder(checkOrder.PaymentCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to get payment", err.Error())
	}

	payment.StatusID = 3
	status := entities.Status{}
	if err := h.Repository.FindStatus(&status, 3); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed get status", err.Error())
		return
	}
	payment.Status = status

	if err := h.Repository.SavePayment(payment); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save udpated payment", err.Error())
	}
	helper.SuccessResponse(c, http.StatusOK, "Update payment Successful", &payment)
}
