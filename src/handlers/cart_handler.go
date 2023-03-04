package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) AddProductToCart(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)
	var newItem entities.CartProduct

	if err := c.BindJSON(&newItem); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// find existing cart
	var cart entities.Cart
	if err := h.Repository.FindCartByUserID(user.ID, &cart); err != nil {
		// create new cart if not exist
		cart.Code = user.ID
		cart.UserID = user.ID

		cartProduct := entities.CartProduct{
			ProductID: uint(newItem.ProductID),
			Quantity:  newItem.Quantity,
		}

		cart.Products = append(cart.Products, cartProduct)

		if err := h.Repository.SaveCart(&cart); err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save cart", nil)
			return
		}
		helper.SuccessResponse(c, http.StatusOK, "Item added to cart", nil)
		return
	}

	// check if product is already in cart
	cartproduct := entities.CartProduct{}
	for _, p := range cart.Products {
		if p.ProductID == uint(newItem.ProductID) {
			cartproduct.Quantity += newItem.Quantity
			if err := h.Repository.SaveCart(&cart); err != nil {
				helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save cart", nil)
				return
			}
			helper.SuccessResponse(c, http.StatusOK, "Item added to cart", nil)
			return
		}
	}

	// add new product to cart
	cartProduct := entities.CartProduct{
		ProductID: uint(newItem.ProductID),
		Quantity:  newItem.Quantity,
	}
	cart.Products = append(cart.Products, cartProduct)

	if err := h.Repository.SaveCart(&cart); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save cart", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Item added to cart", nil)
}

func (h *handler) RemoveItemFromCart(c *gin.Context) {

}

func (h *handler) BuyFromCart(c *gin.Context) {

}

func (h *handler) InstantBuy(c *gin.Context) {

}
