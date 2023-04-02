package services

import (
	"kedaiprogrammer/helpers"
)

type Services interface {
	Save(input AddServiceInput) (Service, error)
	GetAll(filterValue string, size int, page int, field string, dir string, filterField string, filterType string) ([]map[string]interface{}, int, int, error)
	GetService(id string) (map[string]interface{}, error)
}

type services struct {
	repository Repository
}

func NewServices(repository Repository) *services {
	return &services{repository}
}

func (s *services) Save(input AddServiceInput) (Service, error) {
	service := Service{}
	service.ServiceID = helpers.GenerateUUID()
	service.ServiceName = input.ServiceName
	service.IsActive = true
	service.BusinessID = input.BusinessID
	newService, err := s.repository.Save(service)
	if err != nil {
		return service, err
	}
	return newService, nil
}

func (s *services) GetAll(filterValue string, size int, page int, field string, dir string, filterField string, filterType string) ([]map[string]interface{}, int, int, error) {
	return s.repository.GetAllWithCounts(filterValue, size, page, field, dir, filterField, filterType)
}

func (s *services) GetService(id string) (map[string]interface{}, error) {
	return s.repository.GetService(id)
}
