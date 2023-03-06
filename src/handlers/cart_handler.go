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
		cart.UserID = user.ID
		cartProduct := entities.CartProduct{
			ProductID: uint(productID),
			Quantity:  newItem.Quantity,
		}
		h.Repository.CreateCartProduct(&cartProduct)
		cart.Products = append(cart.Products, cartProduct)
		if err := h.Repository.CreateCart(&cart); err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to created cart", nil)
			return
		}
		helper.SuccessResponse(c, http.StatusOK, "Item added to cart", nil)
		return
	}

	cartProduct := entities.CartProduct{
		ProductID: uint(productID),
		Quantity:  newItem.Quantity,
	}

	for i, p := range cart.Products {
		if p.ProductID == cartProduct.ProductID {
			cart.Products[i].Quantity += cartProduct.Quantity
			if err := h.Repository.SaveCart(&cart); err != nil {
				helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save cart", nil)
				return
			}
			helper.SuccessResponse(c, http.StatusOK, "Item added to cart", nil)
			return
		}
	}

	cart.Products = append(cart.Products, cartProduct)
	if err := h.Repository.SaveCart(&cart); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save cart", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product added to cart", nil)
}

func (h *handler) RemoveProductFromCart(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)

	cartFound, err := h.Repository.GetCart(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadGateway, "Can't find the cart", nil)
	}

	productIDStr := c.Query("Product_ID")

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to parse string into uint64", nil)
		return
	}

	if err := h.Repository.DeleteCartProduct(cartFound, uint(productID)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove product from cart", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product deleted from cart", nil)
}

func (h *handler) GetCart(c *gin.Context) {

}
