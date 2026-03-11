package repository

import (
	"context"

	"github.com/uthmanduro/BracketForge/internal/model"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	CreateOrganization(name string) (uint, error)
	GetOrganizationByID(id uint) (*model.Organization, error)
	UpdateOrganization(org *model.Organization) error
	DeleteOrganization(id uint) error
}

type OrganizationRepo struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepo {
	db.AutoMigrate(&model.Organization{}) // Ensure the Organization table is created
	return &OrganizationRepo{db: db}
}

func (r *OrganizationRepo) CreateOrganization(ctx context.Context, name string) (uint, error) {
	org := &model.Organization{Name: name}
	if err := r.db.Create(org).Error; err != nil {
		return 0, err
	}
	return org.ID, nil
}

func (r *OrganizationRepo) GetOrganizationByID(ctx context.Context, id uint) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.First(&org, id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepo) UpdateOrganization(ctx context.Context, org *model.Organization) error {
	return r.db.Save(org).Error
}

func (r *OrganizationRepo) DeleteOrganization(ctx context.Context, id uint) error {
	return r.db.Delete(&model.Organization{}, id).Error
}