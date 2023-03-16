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

	newItem := model.ProductQuantity{}
	if err := c.BindJSON(&newItem); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	productIDStr := c.Query("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to parse string into uint64", nil)
		return
	}

	product, err := h.Repository.GetProductByID(uint(productID))
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Product not found!!!", err.Error())
		return
	}

	// find existing cart
	cart := entities.Cart{}
	if err := h.Repository.GetOrCreateCart(user.ID, &cart); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't find the cart", nil)
		return
	}

	newCartProduct := entities.CartProduct{
		ProductID:    uint(productID),
		Product:      *product,
		Quantity:     newItem.Quantity,
		SellerID:     product.SellerID,
		ProductPrice: float64(newItem.Quantity) * product.Price,
	}

	found := true
	for i, p := range cart.CartProducts {
		if p.ProductID == uint(productID) {
			h.Repository.DeleteCartProductByID(p.ID)
			cart.CartProducts[i] = newCartProduct
			found = false
			break
		}
	}

	if found {
		cart.CartProducts = append(cart.CartProducts, newCartProduct)
	}

	if err := h.Repository.SaveCart(&cart); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save into cart", err.Error())
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

	// if err := h.Repository.SaveCart(cartFound); err != nil {
	// 	helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save into cart", nil)
	// 	return
	// }

	if err := h.Repository.DeleteCartProduct(cartFound.ID, uint(productID)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove product from cart", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product deleted from cart", nil)
}

func (h *handler) GetProductCart(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)
	cartFound, err := h.Repository.GetProductCart(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadGateway, "Can't find the cart", nil)
	}
	helper.SuccessResponse(c, http.StatusOK, "Product find!!!", cartFound)
}
