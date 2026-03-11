package service

import (
	"context"

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

func (s *OrganizationService) CreateOrganization(ctx context.Context, name string) (uint, error) {
	return s.organizationRepo.CreateOrganization(ctx,name)
}

func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id uint) (*model.Organization, error) {
	return s.organizationRepo.GetOrganizationByID(ctx, id)
}

func (s *OrganizationService) UpdateOrganization(ctx context.Context, org *model.Organization) error {
	return s.organizationRepo.UpdateOrganization(ctx, org)
}

func (s *OrganizationService) DeleteOrganization(ctx context.Context, id uint) error {
	return s.organizationRepo.DeleteOrganization(ctx, id)
}