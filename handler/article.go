package handler

import (
	"fmt"
	"kedaiprogrammer/authorization"
	"kedaiprogrammer/helper"
	"kedaiprogrammer/helpers"
	"kedaiprogrammer/master/articles"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
)

type articleHandler struct {
	services articles.Services
}

func NewArticleHandler(services articles.Services) *articleHandler {
	return &articleHandler{services}
}

func (h *articleHandler) CreateData(c *gin.Context) {
	s3 := c.MustGet("S3").(*session.Session)

	current := c.MustGet("current").(*authorization.JWTClaim)
	uuid, err := helpers.Decrypt(current.Uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var request articles.CreateRequest
	err = c.ShouldBind(&request)
	fmt.Println(err)
	if err != nil {
		response := helpers.APIResponse("Create New Category Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newArticle, err := h.services.Save(s3, request, uuid)
	fmt.Println(err)

	if err != nil {
		response := helper.APIResponse("Create New Article Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Create Article Success", http.StatusOK, "success", newArticle)
	c.JSON(http.StatusOK, response)
}

func (h *articleHandler) GetAll(c *gin.Context) {
	search := c.Query("search")
	tag := c.Query("tag")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "1"))
	orderColumn := c.DefaultQuery("order_column", "category_name")
	orderDirection := c.DefaultQuery("order_direction", "asc")

	data, countFiltered, countAll, err := h.services.GetAll(tag, search, limit, offset, orderColumn, orderDirection)
	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusInternalServerError, "success", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	lastPage := countAll / limit
	if countFiltered < countAll {
		lastPage = countFiltered/limit + 1
	}
	if countFiltered < limit {
		lastPage = 1
	}
	response := helpers.APIDTResponse("Success to Get Articles", http.StatusOK, "success", data, countFiltered, countAll, lastPage)
	c.JSON(http.StatusOK, response)
}

func (h *articleHandler) GetDetailArticle(c *gin.Context) {
	article_id := c.Params.ByName("article_id")
	article, err := h.services.GetOne(article_id)
	if err != nil {
		response := helpers.APIResponse(err.Error(), http.StatusInternalServerError, "success", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helpers.APIResponse("Success to Get Article Detail", http.StatusOK, "success", article)
	c.JSON(http.StatusOK, response)
}
