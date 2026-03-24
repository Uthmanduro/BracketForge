package repository

import (

	"github.com/uthmanduro/BracketForge/internal/model"
	"gorm.io/gorm"
)

type PlayerRepository struct {
	DB *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{
		DB: db,
	}
}

func (r *PlayerRepository) Create(player *model.Player) error {
	return r.DB.Create(player).Error
}

func (r *PlayerRepository) GetByID(id, orgID string) (*model.Player, error) {
	var p model.Player
	return &p, r.DB.First(&p, "id = ? AND organization_id = ?", id, orgID).Error
}
 
func (r *PlayerRepository) ListByOrg(orgID string) ([]*model.Player, error) {
	var players []*model.Player
	return players, r.DB.Where("organization_id = ?", orgID).Order("name").Find(&players).Error
}
 
func (r *PlayerRepository) Update(p *model.Player) error {
	return r.DB.Model(p).Updates(map[string]interface{}{
		"name":     p.Name,
		"email":    p.Email,
		"metadata": p.Metadata,
	}).Error
}
 
func (r *PlayerRepository) Delete(id, orgID string) error {
	return r.DB.Delete(&model.Player{}, "id = ? AND organization_id = ?", id, orgID).Error
}