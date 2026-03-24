package service

import (
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

// func (s *OrganizationService) CreateOrganization(ctx context.Context, name string) (uint, error) {
// 	return s.organizationRepo.CreateOrganization(ctx,name)
// }

// func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id uint) (*model.Organization, error) {
// 	return s.organizationRepo.GetOrganizationByID(ctx, id)
// }

func (s *OrganizationService) Create(name string) (*model.Organization, error) {
	org := &model.Organization{Name: name}
	return org, s.organizationRepo.Create(org)
}
 
func (s *OrganizationService) GetByID(id string) (*model.Organization, error) {
	return s.organizationRepo.GetByID(id)
}
