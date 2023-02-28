package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *handler) UserRegister(c *gin.Context) {
	newUser := model.RegisterUser{}
	if err := c.ShouldBindJSON(&newUser); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed create new user!", nil)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 20)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create password", nil)
		return
	}

	user := entities.User{
		FName:    newUser.FName,
		Email:    newUser.Email,
		Username: newUser.Username,
		Password: string(hashPassword),
		Address:  newUser.Address,
		Contact:  newUser.Contact,
		Gender:   newUser.Gender,
	}

	err = h.Repository.CreateUser(&user)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed save new user", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Register Successful", user)

}

func (h *handler) UserLogin(c *gin.Context) {
	loginUser := model.LoginUser{}
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the user", nil)
		return
	}

	user, err := h.Repository.FindUser(&loginUser)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't find the user", nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Can't compare the password", nil)
		return
	}

	//JWT TOKEN
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	//GET JWT TOKEN
	tokenString, err := token.SignedString([]byte(user.Password))

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create token", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login Successful", nil)
	c.JSON(http.StatusOK, gin.H{
		"token JWT": tokenString,
	})
}

func (h *handler) UserUpdate(c *gin.Context) {
	idReq := model.GetUserByID{}
	updateUser := model.UpdateUser{}

	if err := c.ShouldBindUri(&idReq); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Update failed", nil)
		return
	}

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Update failed", nil)
		return
	}

	userFound, err := h.Repository.FindUserByID(idReq.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to find user", nil)
		return
	}

	if updateUser.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 20)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create password", nil)
			return
		}
		userFound.Password = string(hashPassword)
	} else if updateUser.FName != "" {
		userFound.FName = updateUser.FName
	} else if updateUser.Email != "" {
		userFound.Email = updateUser.Email
	} else if updateUser.Username != "" {
		userFound.Username = updateUser.Username
	} else if updateUser.Address != "" {
		userFound.Address = updateUser.Address
	} else if updateUser.Contact != "" {
		userFound.Contact = updateUser.Contact
	} else if updateUser.Gender != userFound.Gender {
		userFound.FName = updateUser.FName
	}

	if err := h.Repository.UpdateUser(userFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed updated user", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &userFound)
}
