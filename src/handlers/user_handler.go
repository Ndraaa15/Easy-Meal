package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *handler) UserRegister(c *gin.Context) {
	newUser := model.RegisterUser{}
	if err := c.ShouldBindJSON(&newUser); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Password format is incorrect. Please follow the specified format and try again!", nil)
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

	if err = h.Repository.CreateUser(&user); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to add user to database. Please try again later or contact customer service for help", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Register successful! Welcome, "+user.Username+"!", user)

}

func (h *handler) UserLogin(c *gin.Context) {
	loginUser := model.LoginUser{}
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	userFound, err := h.Repository.FindUser(&loginUser)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "User not found. Please try again with a valid username!", nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginUser.Password)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Wrong password. Please try again with a valid password!", nil)
		return
	}

	//JWT TOKEN
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      userFound.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"fname":    userFound.FName,
		"id":       userFound.ID,
		"username": userFound.Username,
	})

	//GET JWT TOKEN
	tokenString, err := token.SignedString([]byte(h.config.GetEnv("SECRET_KEY")))

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create token JWT. Please try again to login!", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login successful! Welcome back, "+userFound.Username+"!", nil)
	c.JSON(http.StatusOK, gin.H{
		"JWT Token": tokenString,
	})
}

func (h *handler) UserUpdate(c *gin.Context) {

	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)
	updateUser := model.UpdateUser{}

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", nil)
		return
	}

	userFound, err := h.Repository.FindUserByID(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "User not found. Please try again later!", nil)
		return
	}

	fmt.Println(user.ID)

	if updateUser.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 12)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create new password", nil)
			return
		}
		userFound.Password = string(hashPassword)
	}
	if updateUser.FName != "" {
		userFound.FName = updateUser.FName
	}
	if updateUser.Email != "" {
		userFound.Email = updateUser.Email
	}
	if updateUser.Username != "" {
		userFound.Username = updateUser.Username
	}
	if updateUser.Address != "" {
		userFound.Address = updateUser.Address
	}
	if updateUser.Contact != "" {
		userFound.Contact = updateUser.Contact
	}
	if updateUser.Gender != userFound.Gender {
		userFound.Gender = updateUser.Gender
	}

	if err := h.Repository.UpdateUser(userFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &userFound)
}
