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
	adminClaims, _ := c.Get("admin")
	admin := adminClaims.(model.AdminClaims)

	newProduct := model.NewProduct{}
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed create new product post!", nil)
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
	}

	if err := c.ShouldBindJSON(&UpdateProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	productFound, err := h.Repository.GetProductByID(idProduct.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Product not found", nil)
	}

	if productFound.Name != "" {
		productFound.Name = UpdateProduct.Name
	} else if productFound.Price != 0 {
		productFound.Price = UpdateProduct.Price
	} else if productFound.Stock != 0 {
		productFound.Stock = UpdateProduct.Stock
	} else if productFound.Description != "" {
		productFound.Description = UpdateProduct.Description
	}

	if err := h.Repository.SaveProduct(productFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed update product", nil)
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &productFound)
}

func (h *handler) GetAllProduct(c *gin.Context) {
	products, err := h.Repository.GetAllProduct()
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't find the product", nil)
	}
	helper.SuccessResponse(c, http.StatusBadRequest, "Product found", products)
}

func (h *handler) GetProductByID(c *gin.Context) {
	idProduct := model.GetProductByID{}
	if err := c.ShouldBindUri(&idProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't find id product", nil)
		return
	}
	product, err := h.Repository.GetProductByID(idProduct.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Cant find the product", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Find product successful", product)

}
