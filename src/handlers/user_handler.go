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
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed add user to database. Please try again later!", err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Register successful! Welcome, "+user.Username+"!", &user)
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

	helper.SuccessResponse(c, http.StatusOK, tokenString, &userFound)
}

func (h *handler) UserUpdate(c *gin.Context) {
	userClaims, exist := c.Get("user")
	if !exist {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to load JWT token, please try again!", nil)
		return
	}
	user := userClaims.(model.UserClaims)

	fname := c.PostForm("fname")
	email := c.PostForm("email")
	password := c.PostForm("password")
	username := c.PostForm("username")
	address := c.PostForm("address")
	contact := c.PostForm("contact")
	gender := c.PostForm("gender")

	supClient := supabasestorageuploader.NewSupabaseClient(
		h.config.GetEnv("SUPABASE_URL"),
		h.config.GetEnv("SUPABASE_API_KEY"),
		h.config.GetEnv("SUPABASE_STORAGE"),
		h.config.GetEnv("SUPABASE_USER_FOLDER"),
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
	if password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create new password", err.Error())
			return
		}
		userFound.Password = string(hashPassword)
	}

	if err := h.Repository.UpdateUser(userFound); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to save update to database. Please try again later!", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Update Successful", &userFound)
}
