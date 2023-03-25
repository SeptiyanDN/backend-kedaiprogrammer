package categories

import (
	"kedaiprogrammer/helpers"
)

type Services interface {
	SaveCategory(input AddCategoryInput) (Category, error)
	GetAll(search string, limit int, offset int, OrderColumn string, orderDirection string) ([]map[string]interface{}, int, int, error)
	GetCategory(id string) (map[string]interface{}, error)
}

type services struct {
	repository Repository
}

func NewServices(repository Repository) *services {
	return &services{repository}
}

func (s *services) SaveCategory(input AddCategoryInput) (Category, error) {
	category := Category{}
	category.CategoryID = helpers.GenerateUUID()
	category.CategoryName = input.Category_name
	category.Slug = input.Slug
	category.IsActive = true
	category.BusinessID = input.BusinessID
	newCategory, err := s.repository.Save(category)
	if err != nil {
		return category, err
	}
	return newCategory, nil
}

func (s *services) GetAll(search string, limit int, offset int, OrderColumn string, orderDirection string) ([]map[string]interface{}, int, int, error) {
	return s.repository.GetAllWithCounts(search, limit, offset, OrderColumn, orderDirection)
}

func (s *services) GetCategory(id string) (map[string]interface{}, error) {
	return s.repository.GetCategory(id)
}
