package handlers

import (
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetHistory(c *gin.Context) {
	userClaims, exist := c.Get("user")
	if !exist {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to load JWT token, please try again!", nil)
		return
	}
	user := userClaims.(model.UserClaims)

	histories, err := h.Repository.GetHistory(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Can't load the history, please try again!!!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "History found!!!", &histories)
}
