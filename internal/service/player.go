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
 
func (s *PlayerService) Update(id, orgID string, req *model.UpdatePlayerRequest) (*model.Player, error) {
	p, err := s.PlayerRepo.GetByID(id, orgID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Email != nil {
		p.Email = req.Email
	}
	if req.Metadata != nil {
		p.Metadata = req.Metadata
	}
	return p, s.PlayerRepo.Update(p)
}
 
func (s *PlayerService) Delete(id, orgID string) error {
	return s.PlayerRepo.Delete(id, orgID)
}

func (s *PlayerService) GetByEmail(email, orgID string) (*model.Player, error) {
	return s.PlayerRepo.GetByEmail(email, orgID)
}