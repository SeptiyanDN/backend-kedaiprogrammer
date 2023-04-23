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

	h.serviceServices.Save(input)

	c.JSON(http.StatusOK, gin.H{
		"message": "Service saved successfully",
	})
}

func (h *serviceHandler) GetAllServices(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", "5"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	field := c.DefaultQuery("sorters[0][field]", "service_name")
	dir := c.DefaultQuery("sorters[0][dir]", "asc")
	filterField := c.Query("filters[0][field]") // filter column
	filterType := c.Query("filters[0][type]")   // filter type ( like )
	filterValue := c.Query("filters[0][value]") // search

	data, countAll, countFiltered, err := h.serviceServices.GetAll(filterValue, size, page, field, dir, filterField, filterType)
	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusInternalServerError, "success", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	lastPage := countAll / size
	if countFiltered < countAll {
		lastPage = countFiltered/size + 1
	}
	if countFiltered < size {
		lastPage = 1
	}
	response := helpers.APIDTResponse("Success to Get Services", http.StatusOK, "success", data, countFiltered, countAll, lastPage)
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
