package businesses

import "kedaiprogrammer/helpers"

type Services interface {
	SaveBusiness(input AddBusinessInput) (Business, error)
	FindAll() ([]Business, error)
}

type services struct {
	repository Repository
}

func NewServices(repository Repository) *services {
	return &services{repository}
}

func (s *services) SaveBusiness(input AddBusinessInput) (Business, error) {
	business := Business{}
	business.BusinessID = helpers.GenerateUUID()
	business.Business_name = input.Business_name
	business.Domain = input.Domain
	business.IsActive = true
	newBusiness, err := s.repository.Save(business)
	if err != nil {
		return business, err
	}
	return newBusiness, nil
}

func (s *services) FindAll() ([]Business, error) {
	Business, err := s.repository.FindAll()
	if err != nil {
		return Business, err
	}
	return Business, nil
}
