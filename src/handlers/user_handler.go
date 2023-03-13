package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"fmt"
	"net/http"
	"time"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// type UserHandler interface {
// 	UserRegister(c *gin.Context)
// 	UserLogin(c *gin.Context)
// 	UserUpdate(c *gin.Context)
// 	UserUpdatePassword(c *gin.Context)
// }

// type userHandler struct {
// 	Repository repository.Repository
// 	conf       config.Initializer
// }

// func newUserHandler(repo repository.Repository, conf config.Initializer) *userHandler {
// 	return &userHandler{
// 		Repository: repo,
// 		conf:       conf,
// 	}
// }

func (h *handler) UserRegister(c *gin.Context) {
	newUser := model.RegisterUser{}
	if err := c.ShouldBindJSON(&newUser); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", err.Error())
		fmt.Println(err.Error())
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Password format is incorrect. Please follow the specified format and try again!", err.Error())
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
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed add user to database. Please try again later or contact customer service for help", err.Error())
		fmt.Println(err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Register successful! Welcome, "+user.Username+" ! ", user)
}

func (h *handler) UserLogin(c *gin.Context) {
	loginUser := model.LoginUser{}
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", err.Error())
		return
	}

	userFound, err := h.Repository.FindUser(&loginUser)

	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "User not found. Please try again with a valid username!", err.Error())
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginUser.Password)); err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Wrong password. Please try again with a valid password!", err.Error())
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      userFound.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"id":       userFound.ID,
		"username": userFound.Username,
	})

	tokenString, err := token.SignedString([]byte(h.config.GetEnv("SECRET_KEY")))

	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create token JWT. Please try again to login!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login successful! Welcome back, "+userFound.Username+" ! ", &userFound)
	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

func (h *handler) UserUpdate(c *gin.Context) {

	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)

	fname := c.PostForm("fname")
	email := c.PostForm("email")
	username := c.PostForm("username")
	address := c.PostForm("address")
	contact := c.PostForm("contact")
	gender := c.PostForm("gender")

	supClient := supabasestorageuploader.NewSupabaseClient(
		"https://arcudskzafkijqukfool.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFyY3Vkc2t6YWZraWpxdWtmb29sIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Nzc2NDk3MjksImV4cCI6MTk5MzIyNTcyOX0.CjOVpoFAdq3U-AeAzsuyV6IGcqx2ZnaXjneTis5qd6w",
		"bcc-project",
		"user-image",
	)
	file, err := c.FormFile("user_image")
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to get user image", err.Error())
		return
	}
	link, err := supClient.Upload(file)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload user image", err.Error())
		return
	}

	userFound, err := h.Repository.FindUserByID(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "User not found. Please try again later!", err.Error())
		return
	}

	if fname != "" {
		userFound.FName = fname
	}
	if email != "" {
		userFound.Email = email
	}
	if username != "" {
		userFound.Username = username
	}
	if address != "" {
		userFound.Address = address
	}
	if contact != "" {
		userFound.Contact = contact
	}
	if gender != "" {
		userFound.Gender = gender
	}
	if userFound.UserImage != link {
		userFound.UserImage = link
	}

	if err := h.Repository.UpdateUser(userFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &userFound)
}

func (h *handler) UserUpdatePassword(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)

	userFound, err := h.Repository.FindUserByID(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "User not found. Please try again later!", err.Error())
		return
	}

	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	if err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(oldPassword)); err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Wrong password. Please try again with a valid password!", err.Error())
		return
	}

	if newPassword != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create new password", err.Error())
			return
		}
		userFound.Password = string(hashPassword)
	}

	if err := h.Repository.UpdateUser(userFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Password Successful", &userFound)
}
