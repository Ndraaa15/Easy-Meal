package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

func (h *handler) AddProductToCart(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)

	newItem := model.NewItem{}
	if err := c.BindJSON(&newItem); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	productIDStr := c.Query("Product_ID")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to parse string into uint64", nil)
		return
	}

	// find existing cart
	cart := entities.Cart{}
	if err := h.Repository.FindCartByUserID(user.ID, &cart); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't find the cart", nil)
		return
	}

	cartFound, err := h.Repository.GetCart(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	cartProduct := entities.CartProduct{
		ProductID: uint(productID),
		Quantity:  newItem.Quantity,
	}

	found := false
	for i, p := range cartFound.CartProducts {
		if p.ProductID == uint(productID) {
			cartFound.CartProducts[i].Quantity += cartProduct.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.CartProducts = append(cartFound.CartProducts, cartProduct)
	}

	cart.CartProducts = append(cart.CartProducts, cartProduct)
	if err := h.Repository.SaveCart(&cart); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Product added to cart", nil)
}

func (h *handler) RemoveProductFromCart(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)

	cartFound, err := h.Repository.GetProductCart(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadGateway, "Can't find the cart", nil)
	}

	productIDStr := c.Query("Product_ID")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to parse string into uint64", nil)
		return
	}

	if err := h.Repository.DeleteCartProduct(cartFound.ID, uint(productID)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove product from cart", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product deleted from cart", nil)
}

func (h *handler) GetCart(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)
	cartFound, err := h.Repository.GetProductCart(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadGateway, "Can't find the cart", nil)
	}
	helper.SuccessResponse(c, http.StatusOK, "Product find!!!", cartFound)
}
