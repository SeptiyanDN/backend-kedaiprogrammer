package helper

import (
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}

func Paginate(page, size int, total int64, data interface{}, c *gin.Context) interface{} {
	lastPage := int(math.Ceil(float64(total) / float64(size)))
	prevPage := ""
	nextPage := ""

	if page > 1 {
		prevPage = fmt.Sprintf("%s?page=%d&size=%d", c.Request.URL.Path, page-1, size)
	}

	if page < lastPage {
		nextPage = fmt.Sprintf("%s?page=%d&size=%d", c.Request.URL.Path, page+1, size)
	}

	return gin.H{
		"data":      data,
		"page":      page,
		"size":      size,
		"total":     total,
		"last_page": lastPage,
		"prev_page": prevPage,
		"next_page": nextPage,
	}
}
