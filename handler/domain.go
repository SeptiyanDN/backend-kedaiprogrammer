package handler

import (
	"encoding/json"
	"kedaiprogrammer/domain"
	"kedaiprogrammer/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAllDomainsHandler(c *gin.Context) {
	service := domain.DomainServices{}
	body, err := service.GetDomains()
	if err != nil {
		response := helper.APIResponse("Failed to Get List Domain", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Get List Domain Success", http.StatusOK, "success", body)
	c.JSON(http.StatusOK, response)
}

func GetAvailabiltyDomain(c *gin.Context) {
	service := domain.DomainServices{}
	keyword := c.Params.ByName("keyword")
	body, err := service.GetAvailabiltyDomain(keyword)
	if err != nil {
		response := helper.APIResponse("Failed to Check Domain Avaibility", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Check Domain Avaibility Successfully", http.StatusOK, "success", body)
	c.JSON(http.StatusOK, response)
}

func GetDetailManageDomain(c *gin.Context) {
	service := domain.DomainServices{}
	domain := c.Params.ByName("domain")
	body, err := service.GetDetailManageDomain(domain)
	if err != nil {
		response := helper.APIResponse("Failed to Get Detail Domain ", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Get Detail Manage Domain Successfully", http.StatusOK, "success", body)
	c.JSON(http.StatusOK, response)
}

func GetBalanceAccount(c *gin.Context) {
	service := domain.DomainServices{}
	body := service.GetBalanceAccount()

	response := helper.APIResponse("Get Balance on Provider Domain Successfully", http.StatusOK, "success", body)
	c.JSON(http.StatusOK, response)
}

func GetPriceDomain(c *gin.Context) {
	service := domain.DomainServices{}
	body := service.GetPriceDomain()

	response := helper.APIResponse("Get Price on Provider Domain Successfully", http.StatusOK, "success", body)
	c.JSON(http.StatusOK, response)
}
func NewPrice(c *gin.Context) {
	service := domain.DomainServices{}
	keyword := c.Params.ByName("keyword")
	object1, err := service.GetAvailabiltyDomain(keyword)
	if err != nil {
		response := helper.APIResponse("Failed to Check Domain Avaibility", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	object2 := service.GetPriceDomain()

	// Menggabungkan kedua object ke dalam array
	var result []map[string]interface{}
	for k, v := range object1 {
		tld := strings.Split(k, ".")[1]
		if props, ok := object2[tld]; ok {
			d := map[string]interface{}{
				"domain":       k,
				"avaliability": v,
			}

			for propName, propValue := range props.(map[string]interface{}) {
				d[propName] = propValue
			}

			result = append(result, d)
		}
	}

	// Mengkonversi hasil gabungan menjadi JSON
	jsonResult, _ := json.MarshalIndent(result, "", "  ")

	var res []map[string]interface{}
	json.Unmarshal(jsonResult, &res)

	response := helper.APIResponse("Get Detail Manage Domain Successfully", http.StatusOK, "success", res)
	c.JSON(http.StatusOK, response)
}
