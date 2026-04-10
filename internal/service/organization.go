package service

import (
	"fmt"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
)

type OrganizationService struct {
	// Add any dependencies needed for organization management, e.g. organization repository, etc.
	organizationRepo *repository.OrganizationRepo
}

func NewOrganizationService(organizationRepo *repository.OrganizationRepo) *OrganizationService {
	return &OrganizationService{
		organizationRepo: organizationRepo,
	}
}

func (s *OrganizationService) Create(name string) (*model.Organization, error) {

	org := &model.Organization{Name: name}
	return org, s.organizationRepo.Create(org)
}
 
func (s *OrganizationService) GetByID(id string) (*model.Organization, error) {
	return s.organizationRepo.GetByID(id)
}

func (s *OrganizationService) GetAll() ([]*model.Organization, error) {
	return s.organizationRepo.GetAll()
}

func (s *OrganizationService) GetByName(name string) (*model.Organization, error) {
	org, err := s.organizationRepo.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}
	return org, nil
}