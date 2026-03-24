// package service

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/uthmanduro/BracketForge/internal/model"
// 	"github.com/uthmanduro/BracketForge/internal/repository"
// )

// type TournamentService struct {
// 	TournamentRepo *repository.TournamentRepo
// }

// func NewTournamentService(tournamentRepo *repository.TournamentRepo) *TournamentService {
// 	return &TournamentService{
// 		TournamentRepo: tournamentRepo,
// 	}
// }

// func (s *TournamentService) CreateTournament(ctx context.Context, name, description, format string, organizationID uint, startDate, endDate *time.Time) (*model.Tournament, error) {
// 	//check if tournament with the same name already exists in the organization
// 	existingTournament, err := s.TournamentRepo.GetTournamentByNameAndOrgID(ctx, name, organizationID)
// 	if err == nil && existingTournament != nil {
// 		return nil, fmt.Errorf("tournament with the same name already exists in the organization")
// 	}

// 	tournament := &model.Tournament{
// 		Name:           name,
// 		Description:    description,
// 		Format:         model.TournamentFormat(format),
// 		OrganizationID: organizationID,
// 		StartDate:      startDate,
// 		EndDate:        endDate,
// 	}
// 	err = s.TournamentRepo.CreateTournament(ctx, tournament)
// 	if err != nil {
// 		fmt.Printf("Error creating tournament: %v\n", err)
// 		return nil, err
// 	}
// 	return tournament, nil
// }

// func (s *TournamentService) GetTournamentByID(ctx context.Context, id uint) (*model.Tournament, error) {
// 	return s.TournamentRepo.GetTournamentByID(ctx, id)
// }

// func (s *TournamentService) UpdateTournament(ctx context.Context, tournament *model.Tournament) error {
// 	return s.TournamentRepo.UpdateTournament(ctx, tournament)
// }

// func (s *TournamentService) DeleteTournament(ctx context.Context, id uint) error {
// 	return s.TournamentRepo.DeleteTournament(ctx, id)
// }


package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/uthmanduro/BracketForge/internal/engine"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
)

type TournamentService struct {
	tournRepo      *repository.TournamentRepo
	stageRepo      *repository.StageRepo
	matchRepo      *repository.MatchRepo
	standingsRepo  *repository.StandingsRepo
	playerRepo     *repository.PlayerRepository
	seEngine       *engine.SingleEliminationEngine
	rrEngine       *engine.RoundRobinEngine
	groupEngine    *engine.GroupEngine
	resultEngine   *engine.ResultEngine
	walkoverEngine *engine.WalkoverEngine
}

func NewTournamentService(
	tr *repository.TournamentRepo,
	sr *repository.StageRepo,
	mr *repository.MatchRepo,
	standRepo *repository.StandingsRepo,
	playerRepo *repository.PlayerRepository,
	se *engine.SingleEliminationEngine,
	rr *engine.RoundRobinEngine,
	ge *engine.GroupEngine,
	re *engine.ResultEngine,
	we *engine.WalkoverEngine,
) *TournamentService {
	return &TournamentService{
		tournRepo: tr, stageRepo: sr, matchRepo: mr,
		standingsRepo: standRepo, playerRepo: playerRepo,
		seEngine: se, rrEngine: rr, groupEngine: ge,
		resultEngine: re, walkoverEngine: we,
	}
}

// ── Tournament CRUD ────────────────────────────────────────────────────────

func (s *TournamentService) Create(orgID string, req *model.CreateTournamentRequest) (*model.Tournament, error) {
	t := &model.Tournament{
		OrganizationID: orgID, Name: req.Name, Format: model.TournamentFormat(req.Format),
		BestOf: req.BestOf, ThirdPlaceMatch: req.ThirdPlaceMatch,
		ScoringRules: req.ScoringRules, StartDate: req.StartDate, EndDate: req.EndDate,
	}
	return t, s.tournRepo.Create(t)
}

func (s *TournamentService) GetByID(id, orgID string) (*model.Tournament, error) {
	return s.tournRepo.GetByID(id, orgID)
}

func (s *TournamentService) List(orgID string) ([]*model.Tournament, error) {
	return s.tournRepo.ListByOrg(orgID)
}

func (s *TournamentService) TransitionStatus(id, orgID, next string) error {
	t, err := s.tournRepo.GetByID(id, orgID)
	if err != nil {
		return err
	}
	allowed := map[string]string{"draft": "registration", "registration": "active", "active": "completed"}
	if allowed[string(t.Status)] != next {
		return fmt.Errorf("cannot transition from %s to %s", t.Status, next)
	}
	return s.tournRepo.UpdateStatus(id, orgID, next)
}

// ── Registrations ──────────────────────────────────────────────────────────

func (s *TournamentService) RegisterPlayer(tournID, orgID string, req *model.RegisterPlayerRequest) (*model.TournamentRegistration, error) {
	t, err := s.tournRepo.GetByID(tournID, orgID)
	if err != nil {
		return nil, err
	}
	if t.Status != "registration" && t.Status != "draft" {
		return nil, errors.New("tournament is not accepting registrations")
	}
	if _, err := s.playerRepo.GetByID(req.PlayerID, orgID); err != nil {
		return nil, errors.New("player not found in this organization")
	}
	reg := &model.TournamentRegistration{TournamentID: tournID, PlayerID: req.PlayerID, Seed: req.Seed}
	return reg, s.tournRepo.Register(reg)
}

func (s *TournamentService) ListRegistrations(tournID, orgID string) ([]*model.TournamentRegistration, error) {
	if _, err := s.tournRepo.GetByID(tournID, orgID); err != nil {
		return nil, err
	}
	return s.tournRepo.ListRegistrations(tournID)
}

func (s *TournamentService) UpdateRegistrationStatus(tournID, orgID, regID, status string) error {
	if _, err := s.tournRepo.GetByID(tournID, orgID); err != nil {
		return err
	}
	reg, err := s.tournRepo.GetRegistration(regID)
	if err != nil {
		return err
	}
	if err := s.tournRepo.UpdateRegistrationStatus(regID, status); err != nil {
		return err
	}
	if status == "withdrawn" || status == "disqualified" {
		t, _ := s.tournRepo.GetByID(tournID, orgID)
		if t.Status == "active" {
			return s.walkoverEngine.ProcessWithdrawal(reg.ID)
		}
	}
	return nil
}

// ── Stages ─────────────────────────────────────────────────────────────────

func (s *TournamentService) CreateStage(tournID, orgID string, req *model.CreateStageRequest) (*model.Stage, error) {
	if _, err := s.tournRepo.GetByID(tournID, orgID); err != nil {
		return nil, err
	}
	stage := &model.Stage{
		TournamentID: tournID, Name: req.Name, Type: model.StageType(req.Type),
		StageOrder: req.StageOrder, AdvanceCount: req.AdvanceCount, BestOf: req.BestOf,
	}
	return stage, s.stageRepo.Create(stage)
}

func (s *TournamentService) ListStages(tournID, orgID string) ([]*model.Stage, error) {
	if _, err := s.tournRepo.GetByID(tournID, orgID); err != nil {
		return nil, err
	}
	return s.stageRepo.ListByTournament(tournID)
}

// ── Draw and bracket ───────────────────────────────────────────────────────

func (s *TournamentService) RunGroupDraw(tournID, orgID, stageID string, req *model.DrawRequest) ([]*model.Group, error) {
	t, err := s.tournRepo.GetByID(tournID, orgID)
	if err != nil {
		return nil, err
	}
	stage, err := s.stageRepo.GetByID(stageID)
	if err != nil {
		return nil, err
	}
	if stage.Type != "group" {
		return nil, errors.New("stage is not a group stage")
	}
	regs, err := s.tournRepo.ListActiveRegistrations(tournID)
	if err != nil {
		return nil, err
	}
	return s.groupEngine.Draw(stage, t, regs, req.NumberOfGroups)
}

func (s *TournamentService) GenerateBracket(tournID, orgID, stageID string) ([]*model.Match, error) {
	t, err := s.tournRepo.GetByID(tournID, orgID)
	if err != nil {
		return nil, err
	}
	stage, err := s.stageRepo.GetByID(stageID)
	if err != nil {
		return nil, err
	}
	regs, err := s.tournRepo.ListActiveRegistrations(tournID)
	if err != nil {
		return nil, err
	}
	switch stage.Type {
	case "knockout":
		return s.seEngine.GenerateBracket(stage, t, regs)
	case "round_robin":
		return s.rrEngine.GenerateMatches(stage, t, regs, nil)
	default:
		return nil, fmt.Errorf("use /draw endpoint for group stages")
	}
}

// ── Match operations ───────────────────────────────────────────────────────

func (s *TournamentService) SubmitResult(ctx context.Context, tournID, orgID, matchID string, req *model.SubmitResultRequest) (*model.MatchDetail, error) {
	if _, err := s.tournRepo.GetByID(tournID, orgID); err != nil {
		return nil, err
	}
	return s.resultEngine.Submit(ctx, matchID, req)
}

func (s *TournamentService) CorrectResult(ctx context.Context, tournID, orgID, matchID string, req *model.SubmitResultRequest) (*model.MatchDetail, error) {
	if _, err := s.tournRepo.GetByID(tournID, orgID); err != nil {
		return nil, err
	}
	return s.resultEngine.CorrectResult(ctx, matchID, req)
}

func (s *TournamentService) GetMatchDetail(matchID string) (*model.MatchDetail, error) {
	m, err := s.matchRepo.GetByID(matchID)
	if err != nil {
		return nil, err
	}
	parts, err := s.matchRepo.GetParticipants(matchID)
	if err != nil {
		return nil, err
	}
	sets, err := s.matchRepo.GetSetScores(matchID)
	if err != nil {
		return nil, err
	}
	return &model.MatchDetail{Match: m, Participants: parts, Sets: sets}, nil
}

func (s *TournamentService) ListStageMatches(stageID string) ([]*model.Match, error) {
	return s.matchRepo.ListByStage(stageID)
}

func (s *TournamentService) ScheduleMatch(matchID string, req *model.ScheduleMatchRequest) error {
	return s.matchRepo.ScheduleMatch(matchID, req.ScheduledAt)
}

// ── Standings ──────────────────────────────────────────────────────────────

func (s *TournamentService) GetStageStandings(stageID string) ([]*model.Standings, error) {
	return s.standingsRepo.ListByStage(stageID)
}

func (s *TournamentService) GetGroupStandings(groupID string) ([]*model.Standings, error) {
	return s.groupEngine.RankedStandings(groupID)
}