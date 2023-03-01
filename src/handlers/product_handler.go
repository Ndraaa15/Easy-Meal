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
	supClient := supabasestorageuploader.NewSupabaseClient(
		"PROJECT_URL",
		"PROJECT_API_KEYS",
		"STORAGE_NAME",
		"STORAGE_FOLDER",
	)

	newProduct := model.NewProduct{}

	if err := c.ShouldBindJSON(&newProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed created new post", nil)
	}

	idReq := model.GetAdminByID{}

	if err := c.ShouldBindUri(&idReq); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Bad request", nil)
	}

	adminFound, err := h.Repository.FindAdminByID(idReq.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the admin", nil)
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(400, gin.H{"data": err.Error()})
		return
	}
	link, err := supClient.Upload(file)
	if err != nil {
		c.JSON(500, gin.H{"data": err.Error()})
		return
	}

	product := entities.Product{
		Name:        newProduct.Name,
		Price:       newProduct.Price,
		Description: newProduct.Description,
		Stock:       newProduct.Stock,
		AdminID:     adminFound.ID,
		ImageLink:   link,
	}

	if err := h.Repository.BindingPostAdmin(adminFound, idReq.ID); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't make the post", nil)
	}

	if err := h.Repository.CreateProduct(&product); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't make the product", nil)
	}

	helper.SuccessResponse(c, http.StatusOK, "Create product Successful", product)
}
