package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"net/http"
	"time"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
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
		// Address:  newUser.Address,
		Contact: newUser.Contact,
		// Gender:   newUser.Gender,
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

	helper.SuccessResponse(c, http.StatusOK, "Login successful! Welcome back, "+userFound.Username+"!", tokenString)
}

func (h *handler) UserUpdate(c *gin.Context) {

	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)

	fname := c.PostForm("fname")
	email := c.PostForm("email")
	username := c.PostForm("username")
	gender := c.PostForm("gender")
	address := c.PostForm("address")
	contact := c.PostForm("contact")

	supClient := supabasestorageuploader.NewSupabaseClient(
		"https://arcudskzafkijqukfool.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFyY3Vkc2t6YWZraWpxdWtmb29sIiwicm9sZSI6ImFub24iLCJpYXQiOjE2Nzc2NDk3MjksImV4cCI6MTk5MzIyNTcyOX0.CjOVpoFAdq3U-AeAzsuyV6IGcqx2ZnaXjneTis5qd6w",
		"bcc-project",
		"user-image",
	)
	file, err := c.FormFile("user_image")
	if err != nil {
		c.JSON(400, gin.H{"data": err.Error()})
		return
	}
	link, err := supClient.Upload(file)
	if err != nil {
		c.JSON(500, gin.H{"data": err.Error()})
		return
	}

	userFound, err := h.Repository.FindUserByID(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "User not found. Please try again later!", nil)
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
		helper.ErrorResponse(c, http.StatusInternalServerError, "User not found. Please try again later!", nil)
		return
	}

	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	if err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(oldPassword)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Wrong password. Please try again with a valid password!", nil)
		return
	}

	if newPassword != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create new password", nil)
			return
		}
		userFound.Password = string(hashPassword)
	}

	if err := h.Repository.UpdateUser(userFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Password Successful", &userFound)
}
