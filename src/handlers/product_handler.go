package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) PostProduct(c *gin.Context) {
	newProduct := model.NewProduct{}

	if err := c.ShouldBindJSON(newProduct); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed created new post", nil)
	}

	product := entities.Product{
		Name:        newProduct.Name,
		Price:       newProduct.Price,
		Description: newProduct.Description,
		Stock:       newProduct.Stock,
	}

	fmt.Println(product)
}
