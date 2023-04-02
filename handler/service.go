package handler

import (
	"kedaiprogrammer/helper"
	"kedaiprogrammer/helpers"
	"kedaiprogrammer/master/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type serviceHandler struct {
	serviceServices services.Services
}

func NewServiceHandler(serviceServices services.Services) *serviceHandler {
	return &serviceHandler{serviceServices}
}

func (h *serviceHandler) SaveService(c *gin.Context) {
	var input services.AddServiceInput
	err := c.ShouldBind(&input)
	if err != nil {
		response := helper.APIResponse("Create New Services Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newService, err := h.serviceServices.Save(input)
	if err != nil {
		response := helper.APIResponse("Create New Services Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Create Service Success", http.StatusOK, "success", newService)
	c.JSON(http.StatusOK, response)
}

func (h *serviceHandler) GetAllServices(c *gin.Context) {
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "1"))
	orderColumn := c.DefaultQuery("order_column", "service_name")
	orderDirection := c.DefaultQuery("order_direction", "asc")

	data, countFiltered, countAll, err := h.serviceServices.GetAll(search, limit, offset, orderColumn, orderDirection)
	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusInternalServerError, "success", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIDTResponse("Success to Get Services", http.StatusOK, "success", data, countFiltered, countAll)
	c.JSON(http.StatusOK, response)
}

func (h *serviceHandler) GetDetailService(c *gin.Context) {
	service_id := c.Params.ByName("id")
	service, err := h.serviceServices.GetService(service_id)
	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusInternalServerError, "success", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helpers.APIResponse("Success to Get Service Detail", http.StatusOK, "success", service)
	c.JSON(http.StatusOK, response)
}
