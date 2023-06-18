package handler

import (
	"fmt"
	"kedaiprogrammer/authorization"
	"kedaiprogrammer/helper"
	"kedaiprogrammer/helpers"
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
	fmt.Println(checkUsername)
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
	_, err = h.userServices.SaveToken(newUser.Uuid, token)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Account Has Been Created", http.StatusOK, "success", true)
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
	cekToken, err := h.userServices.GetUserByUUID(loggedinUser.Uuid)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}
	if cekToken.Token != token {
		_, err = h.userServices.SaveToken(cekToken.Uuid, token)
		if err != nil {
			response := helper.APIResponse("Login Account Failed", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	// formatter := users.FormatUserLogin(token)
	response := helper.APIResponse("Login Successfully", http.StatusOK, "success", loggedinUser)
	c.JSON(http.StatusOK, response)
}
func (h *userHandler) CheckUserLoggedIn(c *gin.Context) {
	current := c.MustGet("current").(*authorization.JWTClaim)
	uuid, err := helpers.Decrypt(current.Uuid)
	if err != nil {
		response := helper.APIResponse("Failed To Login! User UnAuthorized", http.StatusUnauthorized, "success", "Failed")
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	fmt.Println(uuid)
	if uuid == "" {
		response := helper.APIResponse("Failed To Login! User UnAuthorized", http.StatusUnauthorized, "success", "Failed")
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	data, err := h.userServices.CheckLoggedIN(current.Username)
	if err != nil {
		response := helper.APIResponse("Failed To Login! User UnAuthorized", http.StatusUnauthorized, "success", "Failed")
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	response := helper.APIResponse("Success To Check Logged In Status", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
