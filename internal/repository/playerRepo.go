package repository

import (
	"context"

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

func (r *PlayerRepository) CreatePlayer(ctx context.Context, player *model.Player) error {
	return r.DB.Create(player).Error
}

func (r *PlayerRepository) GetPlayerByID(ctx context.Context, id uint) (*model.Player, error) {
	var player model.Player
	if err := r.DB.First(&player, id).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *PlayerRepository) GetPlayersByTournamentID(ctx context.Context, tournamentID uint) ([]model.Player, error) {
	var players []model.Player
	if err := r.DB.Where("tournament_id = ?", tournamentID).Preload("Tournament").Preload("Tournament.Organization").Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}

func (r *PlayerRepository) UpdatePlayer(ctx context.Context, player *model.Player) error {
	return r.DB.Save(player).Error
}

func (r *PlayerRepository) DeletePlayer(ctx context.Context, id uint) error {
	return r.DB.Delete(&model.Player{}, id).Error
}