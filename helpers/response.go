package helpers

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type ResponseDT struct {
	Meta          Meta        `json:"meta"`
	Data          interface{} `json:"data"`
	TotalData     int         `json:"recordsTotal"`
	TotalFiltered int         `json:"recordsFiltered"`
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

	jsonResponseBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		return Response{}
	}

	parsedResponse := map[string]interface{}{}

	err = json.Unmarshal(jsonResponseBytes, &parsedResponse)
	if err != nil {
		return Response{}
	}

	parsedResponse["meta"].(map[string]interface{})["message"] = message

	jsonResponseBytes, err = json.Marshal(parsedResponse)
	if err != nil {
		return Response{}
	}

	err = json.Unmarshal(jsonResponseBytes, &jsonResponse)
	if err != nil {
		return Response{}
	}

	return jsonResponse
}

func APIDTResponse(message string, code int, status string, data interface{}, totalData, totalFiltered int) ResponseDT {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := ResponseDT{
		Meta:          meta,
		Data:          data,
		TotalData:     totalData,
		TotalFiltered: totalFiltered,
	}

	jsonResponseBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		return ResponseDT{}
	}

	parsedResponse := map[string]interface{}{}

	err = json.Unmarshal(jsonResponseBytes, &parsedResponse)
	if err != nil {
		return ResponseDT{}
	}

	parsedResponse["meta"].(map[string]interface{})["message"] = message

	jsonResponseBytes, err = json.Marshal(parsedResponse)
	if err != nil {
		return ResponseDT{}
	}

	err = json.Unmarshal(jsonResponseBytes, &jsonResponse)
	if err != nil {
		return ResponseDT{}
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
