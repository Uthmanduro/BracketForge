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
	return &OrganizationRepo{db: db}
}

func (r *OrganizationRepo) Create(org *model.Organization) error {
	return r.db.Create(org).Error
}

func (r *OrganizationRepo) GetAll() ([]*model.Organization, error) {
	var orgs []*model.Organization
	if err := r.db.Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *OrganizationRepo) GetByID(id string) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepo) GetByName(name string) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.Where("name = ?", name).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepo) UpdateOrganization(ctx context.Context, org *model.Organization) error {
	return r.db.Save(org).Error
}

func (r *OrganizationRepo) DeleteOrganization(ctx context.Context, id string) error {
	return r.db.Delete(&model.Organization{}, id).Error
}