package service

import (
	"context"
	"fmt"
	"time"

	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
)

type TournamentService struct {
	TournamentRepo *repository.TournamentRepo
}

func NewTournamentService(tournamentRepo *repository.TournamentRepo) *TournamentService {
	return &TournamentService{
		TournamentRepo: tournamentRepo,
	}
}

func (s *TournamentService) CreateTournament(ctx context.Context, name, description, format string, organizationID uint, startDate, endDate *time.Time) (*model.Tournament, error) {
	//check if tournament with the same name already exists in the organization
	existingTournament, err := s.TournamentRepo.GetTournamentByNameAndOrgID(ctx, name, organizationID)
	if err == nil && existingTournament != nil {
		return nil, fmt.Errorf("tournament with the same name already exists in the organization")
	}

	tournament := &model.Tournament{
		Name:           name,
		Description:    description,
		Format:         model.TournamentFormat(format),
		OrganizationID: organizationID,
		StartDate:      startDate,
		EndDate:        endDate,
	}
	err = s.TournamentRepo.CreateTournament(ctx, tournament)
	if err != nil {
		fmt.Printf("Error creating tournament: %v\n", err)
		return nil, err
	}
	return tournament, nil
}

func (s *TournamentService) GetTournamentByID(ctx context.Context, id uint) (*model.Tournament, error) {
	return s.TournamentRepo.GetTournamentByID(ctx, id)
}

func (s *TournamentService) UpdateTournament(ctx context.Context, tournament *model.Tournament) error {
	return s.TournamentRepo.UpdateTournament(ctx, tournament)
}

func (s *TournamentService) DeleteTournament(ctx context.Context, id uint) error {
	return s.TournamentRepo.DeleteTournament(ctx, id)
}