package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"

	"net/http"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
)

func (h *handler) PostProduct(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	admin := sellerClaims.(model.SellerClaims)

	newProduct := model.NewProduct{}
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	product := entities.Product{
		Name:        newProduct.Name,
		Price:       newProduct.Price,
		Description: newProduct.Description,
		Stock:       newProduct.Stock,
		AdminID:     admin.ID,
	}

	if err := h.Repository.CreateProduct(&product); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't make the product", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Create product Successful", product)
}

func (h *handler) PostImageProduct(c *gin.Context) {
	supClient := supabasestorageuploader.NewSupabaseClient(
		"https://arcudskzafkijqukfool.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFyY3Vkc2t6YWZraWpxdWtmb29sIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Nzc2NDk3MjksImV4cCI6MTk5MzIyNTcyOX0.CjOVpoFAdq3U-AeAzsuyV6IGcqx2ZnaXjneTis5qd6w",
		"bcc-project",
		"product-image",
	)

	file, err := c.FormFile("product")
	if err != nil {
		c.JSON(400, gin.H{"data": err.Error()})
		return
	}
	link, err := supClient.Upload(file)
	if err != nil {
		c.JSON(500, gin.H{"data": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": link})

}

func (h *handler) UpdateProduct(c *gin.Context) {
	idProduct := model.GetProductByID{}
	UpdateProduct := model.UpdateProduct{}

	if err := c.ShouldBindUri(&idProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}

	if err := c.ShouldBindJSON(&UpdateProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	productFound, err := h.Repository.GetProductByID(idProduct.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Product not found. Please try again later!", nil)
		return
	}

	if productFound.Name != "" {
		productFound.Name = UpdateProduct.Name
	}
	if productFound.Price != 0 {
		productFound.Price = UpdateProduct.Price
	}
	if productFound.Stock != 0 {
		productFound.Stock = UpdateProduct.Stock
	}
	if productFound.Description != "" {
		productFound.Description = UpdateProduct.Description
	}

	if err := h.Repository.SaveProduct(productFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed update product", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &productFound)
}

func (h *handler) GetAllProduct(c *gin.Context) {
	products, err := h.Repository.GetAllProduct()
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Product not found. Please try again later!", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Product found", products)
}

func (h *handler) GetProductByID(c *gin.Context) {
	idProduct := model.GetProductByID{}
	if err := c.ShouldBindUri(&idProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't find the id product", nil)
		return
	}
	product, err := h.Repository.GetProductByID(idProduct.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Product not found. Please try again later!", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Find product successful", product)
}

func (h *handler) DeleteProductByID(c *gin.Context) {
	idProduct := model.GetProductByID{}
	if err := c.ShouldBindUri(&idProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't find the ID product", nil)
		return
	}
	product, err := h.Repository.DeleteProductByID(idProduct.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the product. Please try again!", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Product with name "+product.Name+" succesful deleted", nil)
}
