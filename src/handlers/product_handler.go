package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"strconv"

	"net/http"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
)

// -----------------FOR SELLER----------------------

func (h *handler) PostProduct(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)

	sellerFound, err := h.Repository.FindSellerByID(seller.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Seller not found", err.Error())
	}

	name := c.PostForm("name")
	price := c.PostForm("price")
	description := c.PostForm("description")
	stock := c.PostForm("stock")
	category := c.PostForm("category")
	mass := c.PostForm("mass")

	supClient := supabasestorageuploader.NewSupabaseClient(
		"https://arcudskzafkijqukfool.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFyY3Vkc2t6YWZraWpxdWtmb29sIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Nzc2NDk3MjksImV4cCI6MTk5MzIyNTcyOX0.CjOVpoFAdq3U-AeAzsuyV6IGcqx2ZnaXjneTis5qd6w",
		"bcc-project",
		"product-image",
	)
	file, err := c.FormFile("product_image")
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to get product image", err.Error())
		return
	}
	link, err := supClient.Upload(file)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload product image", err.Error())
		return
	}

	massConv, err := strconv.Atoi(mass)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Error while parsing string into uint", err.Error())
	}

	priceConv, err := strconv.ParseFloat(price, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Error while parsing string into float64", err.Error())
	}

	stockConv, err := strconv.ParseUint(stock, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Error while parsing string into uint", err.Error())
	}

	categoryConv, err := strconv.ParseUint(category, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Error while parsing string into uint", err.Error())
	}

	productCategory := entities.Category{}
	if err := h.Repository.FindCategory(&productCategory, uint(categoryConv)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Sorry, Cant find the category. Please choose another category", err.Error())
		return
	}

	product := entities.Product{
		Name:         name,
		Price:        priceConv,
		Description:  description,
		Stock:        uint(stockConv),
		SellerID:     seller.ID,
		ProductImage: link,
		CategoryID:   uint(categoryConv),
		Category:     productCategory,
		Mass:         uint(massConv),
		Seller:       *sellerFound,
	}

	if err := h.Repository.CreateProduct(&product); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create new product, please try again!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Create product successful!!!", &product)
}

func (h *handler) UpdateProduct(c *gin.Context) {
	productID := model.GetProductByID{}

	if err := c.ShouldBindUri(&productID); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't bind the URI", err.Error())
		return
	}

	name := c.PostForm("name")
	price := c.PostForm("price")
	description := c.PostForm("description")
	stock := c.PostForm("stock")
	category := c.PostForm("category")
	mass := c.PostForm("mass")

	massConv, err := strconv.Atoi(mass)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Error while parsing string into uint", err.Error())
	}

	priceConv, err := strconv.ParseFloat(price, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Error while parsing string into float64", err.Error())
	}

	stockConv, err := strconv.ParseUint(stock, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Error while parsing string into uint", err.Error())
	}

	categoryConv, err := strconv.ParseUint(category, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Error while parsing string into uint", err.Error())
	}

	productCategory := entities.Category{}
	if err := h.Repository.FindCategory(&productCategory, uint(categoryConv)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Sorry, Cant find the category. Please choose another category", err.Error())
		return
	}

	supClient := supabasestorageuploader.NewSupabaseClient(
		"https://arcudskzafkijqukfool.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFyY3Vkc2t6YWZraWpxdWtmb29sIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Nzc2NDk3MjksImV4cCI6MTk5MzIyNTcyOX0.CjOVpoFAdq3U-AeAzsuyV6IGcqx2ZnaXjneTis5qd6w",
		"bcc-project",
		"product-image",
	)
	file, err := c.FormFile("product_image")
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to get product image", err.Error())
		return
	}
	link, err := supClient.Upload(file)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload product image", err.Error())
		return
	}

	productFound, err := h.Repository.GetProductByID(productID.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Product not found. Please try again later!", err.Error())
		return
	}

	if name != "" {
		productFound.Name = name
	}
	if priceConv != 0 {
		productFound.Price = priceConv
	}
	if stockConv != 0 {
		productFound.Stock = uint(stockConv)
	}
	if description != "" {
		productFound.Description = description
	}
	if productFound.ProductImage != link {
		productFound.ProductImage = link
	}
	if mass != "" {
		productFound.Mass = uint(massConv)
	}
	if category != "" {
		productFound.CategoryID = uint(categoryConv)
		productFound.Category = productCategory
	}

	if err := h.Repository.SaveProduct(productFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed update data product", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Product Successful", &productFound)
}

func (h *handler) GetSellerProduct(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)

	products, err := h.Repository.GetSellerProduct(seller.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the product", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product found!!!", &products)
}

func (h *handler) GetSellerProductByID(c *gin.Context) {
	productID := model.GetProductByID{}
	if err := c.ShouldBindUri(&productID); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't bind the URI", nil)
		return
	}

	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)

	productFound, err := h.Repository.GetSellerProductByID(seller.ID, productID.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the product with this seller ID or product ID", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Find product successful", &productFound)
}

func (h *handler) DeleteProductByID(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)

	productID := model.GetProductByID{}
	if err := c.ShouldBindUri(&productID); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't bind the URI", err.Error())
		return
	}

	if err := h.Repository.DeleteProductByID(seller.ID, productID.ID); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the product from database. Please try again!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product succesful deleted", nil)
}

// -----------------FOR BUYER----------------------

func (h *handler) GetAllProduct(c *gin.Context) {
	products, err := h.Repository.GetAllProduct()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Product not found. Please try again later!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product found!!!", &products)
}

func (h *handler) GetProductByID(c *gin.Context) {
	productID := model.GetProductByID{}
	if err := c.ShouldBindUri(&productID); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't bind the URI", err.Error())
		return
	}

	product, err := h.Repository.GetProductByID(productID.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Product not found. Please try again later!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product found!!!", product)
}

func (h *handler) GetProductByFilter(c *gin.Context) {
	categoryIDStr := c.Param("category")
	categoryIDConv, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Error while parse string into uint", err.Error())
	}

	products, err := h.Repository.FilteredProduct(uint(categoryIDConv))
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to get filtered products", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Product filter found", &products)
}

func (h *handler) SearchProduct(c *gin.Context) {
	search := c.Query("search")
	products, err := h.Repository.SearchProduct(search)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Cant find the product in database. Please try again!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Products found!!!", &products)
}
