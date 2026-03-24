package service

import (

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

// func (s *PlayerService) CreatePlayer(ctx context.Context, name string, tournamentID uint, seed, ranking int) (*model.Player, error) {
// 	player := &model.Player{
// 		Name: name,
// 		TournamentID: tournamentID,
// 		Seed: seed,
// 		Ranking: &ranking,
// 	}
// 	err := s.PlayerRepo.CreatePlayer(ctx, player)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return player, nil
// }

// func (s *PlayerService) GetPlayerByID(ctx context.Context, id uint) (*model.Player, error) {
// 	return s.PlayerRepo.GetPlayerByID(ctx, id)
// }

// func (s *PlayerService) GetPlayersByTournamentID(ctx context.Context, tournamentID uint) ([]model.Player, error) {
// 	return s.PlayerRepo.GetPlayersByTournamentID(ctx, tournamentID)
// }

// func (s *PlayerService) UpdatePlayer(ctx context.Context, player *model.Player) error {
// 	return s.PlayerRepo.UpdatePlayer(ctx, player)
// }

// func (s *PlayerService) DeletePlayer(ctx context.Context, id uint) error {
// 	return s.PlayerRepo.DeletePlayer(ctx, id)
// }

func (s *PlayerService) Create(orgID string, req *model.CreatePlayerRequest) (*model.Player, error) {
	p := &model.Player{OrganizationID: orgID, Name: req.Name, Email: req.Email, Metadata: req.Metadata}
	return p, s.PlayerRepo.Create(p)
}
 
func (s *PlayerService) GetByID(id, orgID string) (*model.Player, error) {
	return s.PlayerRepo.GetByID(id, orgID)
}
 
func (s *PlayerService) List(orgID string) ([]*model.Player, error) {
	return s.PlayerRepo.ListByOrg(orgID)
}
 
func (s *PlayerService) Update(id, orgID string, req *model.CreatePlayerRequest) (*model.Player, error) {
	p, err := s.PlayerRepo.GetByID(id, orgID)
	if err != nil {
		return nil, err
	}
	p.Name = req.Name
	p.Email = req.Email
	p.Metadata = req.Metadata
	return p, s.PlayerRepo.Update(p)
}
 
func (s *PlayerService) Delete(id, orgID string) error {
	return s.PlayerRepo.Delete(id, orgID)
}