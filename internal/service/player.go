package service

import (
	"context"

	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
)

type PlayerService struct {
	PlayerRepo *repository.PlayerRepository
}

func NewPlayerService(playerRepo *repository.PlayerRepository) *PlayerService {
	return &PlayerService{
		PlayerRepo: playerRepo,
	}
}

func (s *PlayerService) CreatePlayer(ctx context.Context, name string, tournamentID uint, seed, ranking int) (*model.Player, error) {
	player := &model.Player{
		Name: name,
		TournamentID: tournamentID,
		Seed: seed,
		Ranking: &ranking,
	}
	err := s.PlayerRepo.CreatePlayer(ctx, player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (s *PlayerService) GetPlayerByID(ctx context.Context, id uint) (*model.Player, error) {
	return s.PlayerRepo.GetPlayerByID(ctx, id)
}

func (s *PlayerService) GetPlayersByTournamentID(ctx context.Context, tournamentID uint) ([]model.Player, error) {
	return s.PlayerRepo.GetPlayersByTournamentID(ctx, tournamentID)
}

func (s *PlayerService) UpdatePlayer(ctx context.Context, player *model.Player) error {
	return s.PlayerRepo.UpdatePlayer(ctx, player)
}

func (s *PlayerService) DeletePlayer(ctx context.Context, id uint) error {
	return s.PlayerRepo.DeletePlayer(ctx, id)
}