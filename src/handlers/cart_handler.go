package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) AddProductToCart(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)
	newItem := entities.CartProduct{}
	if err := c.BindJSON(&newItem); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	productID, err := strconv.ParseUint(c.Query("product_id"), 32, 10)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to parse string into uint32", nil)
		return
	}
	// find existing cart
	cart := entities.Cart{}
	if err := h.Repository.FindCartByUserID(user.ID, &cart); err != nil {
		cart.UserID = user.ID
		cartProduct := entities.CartProduct{
			ProductID: uint(productID),
			Quantity:  newItem.Quantity,
		}
		cart.Products = append(cart.Products, cartProduct)
		if err := h.Repository.CreateCart(&cart); err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed save product to cart", nil)
			return
		}
		helper.SuccessResponse(c, http.StatusOK, "Item added to cart", nil)
		return
	}
	cartproduct := entities.CartProduct{}
	for _, p := range cart.Products {
		if p.ProductID == uint(newItem.ProductID) {
			cartproduct.Quantity = newItem.Quantity
			if err := h.Repository.SaveCart(&cart); err != nil {
				helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save cart", nil)
				return
			}
			helper.SuccessResponse(c, http.StatusOK, "Item added to cart", nil)
			return
		}
	}
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

func (h *handler) GetProductFromCart(c *gin.Context) {

}

func (h *handler) RemoveProductFromCart(c *gin.Context) {

}
