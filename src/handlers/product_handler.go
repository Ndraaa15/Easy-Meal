package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"log"
	"strconv"

	"net/http"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
)

func (h *handler) PostProduct(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)

	name := c.PostForm("name")
	price := c.PostForm("price")
	description := c.PostForm("description")
	stock := c.PostForm("stock")
	category := c.PostForm("category")

	supClient := supabasestorageuploader.NewSupabaseClient(
		"https://arcudskzafkijqukfool.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFyY3Vkc2t6YWZraWpxdWtmb29sIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Nzc2NDk3MjksImV4cCI6MTk5MzIyNTcyOX0.CjOVpoFAdq3U-AeAzsuyV6IGcqx2ZnaXjneTis5qd6w",
		"bcc-project",
		"product-image",
	)
	file, err := c.FormFile("product_image")
	if err != nil {
		c.JSON(400, gin.H{"data": err.Error()})
		return
	}
	link, err := supClient.Upload(file)
	if err != nil {
		c.JSON(500, gin.H{"data": err.Error()})
		return
	}

	priceConv, _ := strconv.ParseFloat(price, 64)
	stockConv, _ := strconv.ParseUint(stock, 10, 64)
	categoryConv, _ := strconv.ParseUint(category, 10, 64)

	categoryProduct := entities.Category{}
	if err := h.Repository.FindCategory(&categoryProduct, uint(categoryConv)); err != nil {
		helper.ErrorResponse(c, http.StatusBadGateway, err.Error(), nil)
		return
	}
	log.Println(categoryProduct)
	product := entities.Product{
		Name:         name,
		Price:        priceConv,
		Description:  description,
		Stock:        uint(stockConv),
		SellerID:     seller.ID,
		ProductImage: link,
		CategoryID:   uint(categoryConv),
		Category:     categoryProduct,
	}

	if err := h.Repository.CreateProduct(&product); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't make the product", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Create product Successful", &product)
}

func (h *handler) UpdateProduct(c *gin.Context) {
	idProduct := model.GetProductByID{}

	if err := c.ShouldBindUri(&idProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}

	name := c.PostForm("name")
	price := c.PostForm("price")
	description := c.PostForm("description")
	stock := c.PostForm("stock")
	// category := c.PostForm("category")

	priceConv, _ := strconv.ParseFloat(price, 64)
	stockConv, _ := strconv.ParseUint(stock, 10, 64)

	supClient := supabasestorageuploader.NewSupabaseClient(
		"https://arcudskzafkijqukfool.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFyY3Vkc2t6YWZraWpxdWtmb29sIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Nzc2NDk3MjksImV4cCI6MTk5MzIyNTcyOX0.CjOVpoFAdq3U-AeAzsuyV6IGcqx2ZnaXjneTis5qd6w",
		"bcc-project",
		"product-image",
	)
	file, err := c.FormFile("product_image")
	if err != nil {
		c.JSON(400, gin.H{"data": err.Error()})
		return
	}
	link, err := supClient.Upload(file)
	if err != nil {
		c.JSON(500, gin.H{"data": err.Error()})
		return
	}

	productFound, err := h.Repository.GetProductByID(idProduct.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Product not found. Please try again later!", nil)
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
	// if category != "" {
	// 	productFound.Category = category
	// }

	if err := h.Repository.SaveProduct(productFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed update product", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &productFound)
}

func (h *handler) GetSellerProduct(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)
	products, err := h.Repository.GetSellerProduct(seller.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the product", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Product Find!!!", products)
}

func (h *handler) GetAllProduct(c *gin.Context) {
	pageStr := c.Param("page")
	page, _ := strconv.Atoi(pageStr)

	offset := (page - 1) * 12

	products, err := h.Repository.GetAllProduct(uint(offset))
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

func (h *handler) GetSellerProductByID(c *gin.Context) {
	idProduct := model.GetProductByID{}
	if err := c.ShouldBindUri(&idProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't find the id product", nil)
		return
	}
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)
	productFound, err := h.Repository.GetSellerProductByID(seller.ID, idProduct.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't find the product with this seller id or product id", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Find product successful", &productFound)
}

func (h *handler) DeleteProductByID(c *gin.Context) {
	sellerClaims, _ := c.Get("seller")
	seller := sellerClaims.(model.SellerClaims)
	idProduct := model.GetProductByID{}
	if err := c.ShouldBindUri(&idProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Can't find the ID product", nil)
		return
	}
	product, err := h.Repository.DeleteProductByID(seller.ID, idProduct.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete the product. Please try again!", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Product with name "+product.Name+" succesful deleted", nil)
}

func (h *handler) GetProductByFilter(c *gin.Context) {
	categoryIDStr := c.Param("category")
	categoryIDConv, _ := strconv.Atoi(categoryIDStr)
	products, err := h.Repository.FilteredProduct(uint(categoryIDConv))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to filter", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Product filter found", products)

}

func (h *handler) SearchProduct(c *gin.Context) {
	search := c.Query("search")
	product, err := h.Repository.SearchProduct(search)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Cant find the product!!!", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Product found", &product)
}
