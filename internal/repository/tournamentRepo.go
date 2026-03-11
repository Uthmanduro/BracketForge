package repository

import (
	"context"

	"github.com/uthmanduro/BracketForge/internal/model"
	"gorm.io/gorm"
)

type TournamentRepo struct {
	DB *gorm.DB
}

func NewTournamentRepo(db *gorm.DB) *TournamentRepo {
	return &TournamentRepo{
		DB: db,
	}
}

func (r *TournamentRepo) CreateTournament(ctx context.Context, tournament *model.Tournament) error {
	return r.DB.Create(tournament).Error
}

func (r *TournamentRepo) GetTournamentByNameAndOrgID(ctx context.Context, name string, orgID uint) (*model.Tournament, error) {
	var tournament model.Tournament
	err := r.DB.Where("name = ? AND organization_id = ?", name, orgID).First(&tournament).Error
	if err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (r *TournamentRepo) GetTournamentByID(ctx context.Context, id uint) (*model.Tournament, error) {
	var tournament model.Tournament
	err := r.DB.First(&tournament, id).Error
	if err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (r *TournamentRepo) UpdateTournament(ctx context.Context, tournament *model.Tournament) error {
	return r.DB.Save(tournament).Error
}

func (r *TournamentRepo) DeleteTournament(ctx context.Context, id uint) error {
	return r.DB.Delete(&model.Tournament{}, id).Error
}