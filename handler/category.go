package handler

import (
	"kedaiprogrammer/categories"
	"kedaiprogrammer/helper"
	"kedaiprogrammer/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type categoryHandler struct {
	categoryServices categories.Services
}

func NewCategoryHandler(categoryServices categories.Services) *categoryHandler {
	return &categoryHandler{categoryServices}
}

func (h *categoryHandler) SaveCategory(c *gin.Context) {
	var input categories.AddCategoryInput
	err := c.ShouldBind(&input)
	if err != nil {
		response := helper.APIResponse("Create New Category Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newCategory, err := h.categoryServices.SaveCategory(input)
	if err != nil {
		response := helper.APIResponse("Create New Category Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Create Category Success", http.StatusOK, "success", newCategory)
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) GetAllCategory(c *gin.Context) {
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "1"))
	orderColumn := c.DefaultQuery("order_column", "category_name")
	orderDirection := c.DefaultQuery("order_direction", "asc")

	data, countFiltered, countAll, err := h.categoryServices.GetAll(search, limit, offset, orderColumn, orderDirection)
	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusInternalServerError, "success", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIDTResponse("Success to Get Categories", http.StatusOK, "success", data, countFiltered, countAll)
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) GetDetailCategory(c *gin.Context) {
	category_id := c.Params.ByName("id")
	category, err := h.categoryServices.GetCategory(category_id)
	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusInternalServerError, "success", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helpers.APIResponse("Success to Get Category Detail", http.StatusOK, "success", category)
	c.JSON(http.StatusOK, response)
}
