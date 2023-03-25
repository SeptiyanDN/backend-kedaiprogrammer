package handler

import (
	"kedaiprogrammer/authorization"
	"kedaiprogrammer/helper"
	"kedaiprogrammer/users"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandler struct {
	userServices users.Services
	authServices authorization.Services
}

func generateUUID() string {
	u := uuid.New()
	return u.String()
}

func NewUserHandler(userServices users.Services, authServices authorization.Services) *userHandler {
	return &userHandler{userServices, authServices}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input users.RegisterUserInput
	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var inputCheckEmail users.CheckEmailInput
	inputCheckEmail.Email = input.Email
	cekEmail, _ := h.userServices.IsEmailAvailable(inputCheckEmail)

	if cekEmail == false {
		response := helper.APIResponse("Email Sudah Pernah Terdaftar", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputCheckUsername users.CheckUsernameInput
	inputCheckUsername.Username = input.Username
	checkUsername, _ := h.userServices.IsUsernameAvailable(inputCheckUsername)

	if checkUsername == false {
		response := helper.APIResponse("Username Sudah Pernah Terdaftar", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newUser, err := h.userServices.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authServices.GenerateJWT(newUser.Uuid, newUser.Username)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	generateUUID := generateUUID()
	_, err = h.userServices.SaveUUID(newUser.ID, generateUUID)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.userServices.SaveToken(newUser.ID, token)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := users.FormatUser(newUser, token, generateUUID)
	response := helper.APIResponse("Account Has Been Created", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input users.LoginInput
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userServices.Login(input)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authServices.GenerateJWT(loggedinUser.Uuid, loggedinUser.Username)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	cekToken, err := h.userServices.GetUsersByID(loggedinUser.ID)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}
	if cekToken.Token != token {
		_, err = h.userServices.SaveToken(cekToken.ID, token)
		if err != nil {
			response := helper.APIResponse("Login Account Failed", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	formatter := users.FormatUserLogin(token)
	response := helper.APIResponse("Login Successfully", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}
