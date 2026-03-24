package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"

	"gorm.io/gorm"
)

type ResultEngine struct {
	store         *repository.Store
	matchRepo     *repository.MatchRepo
	standingsRepo *repository.StandingsRepo
	stageRepo     *repository.StageRepo
}

func NewResultEngine(
	mr *repository.MatchRepo,
	sr *repository.StandingsRepo,
	stageRepo *repository.StageRepo,
	store *repository.Store,
) *ResultEngine {
	return &ResultEngine{store: store, matchRepo: mr, standingsRepo: sr, stageRepo: stageRepo}
}

// Submit records a full match result atomically.
func (e *ResultEngine) Submit(ctx context.Context, matchID string, req *model.SubmitResultRequest) (*model.MatchDetail, error) {
	// Pre-flight reads — outside the transaction.
	m, err := e.matchRepo.GetByID(matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}
	if m.Status != "pending" && m.Status != "in_progress" {
		return nil, fmt.Errorf("match is %s, cannot submit result", m.Status)
	}
	participants, err := e.matchRepo.GetParticipants(matchID)
	if err != nil {
		return nil, err
	}
	if len(participants) != 2 {
		return nil, fmt.Errorf("match needs 2 participants, has %d", len(participants))
	}
	if err := validateSets(req.Sets, m.BestOf); err != nil {
		return nil, err
	}

	// Tally set and game counts.
	p1Sets, p2Sets, p1Games, p2Games := 0, 0, 0, 0
	for _, s := range req.Sets {
		if s.P1Games > s.P2Games {
			p1Sets++
		} else {
			p2Sets++
		}
		p1Games += s.P1Games
		p2Games += s.P2Games
	}

	setsToWin := (m.BestOf / 2) + 1
	if p1Sets < setsToWin && p2Sets < setsToWin {
		return nil, fmt.Errorf("match not decided: need %d sets to win, p1=%d p2=%d",
			setsToWin, p1Sets, p2Sets)
	}

	// Identify winner and loser by slot.
	winnerPart := participantBySlot(participants, 1)
	loserPart := participantBySlot(participants, 2)
	if p2Sets >= setsToWin {
		winnerPart, loserPart = loserPart, winnerPart
	}

	winnerSets, loserSets := p1Sets, p2Sets
	winnerGames, loserGames := p1Games, p2Games
	if winnerPart.Slot == 2 {
		winnerSets, loserSets = p2Sets, p1Sets
		winnerGames, loserGames = p2Games, p1Games
	}

	stage, err := e.stageRepo.GetByID(m.StageID)
	if err != nil {
		return nil, err
	}

	// All writes in one atomic transaction.
	err = e.store.RunInTx(ctx, func(tx *gorm.DB) error {
		now := time.Now()

		// 1. Insert set scores.
		for _, s := range req.Sets {
			ss := &model.SetScore{
				MatchID:          matchID,
				SetNumber:        s.SetNumber,
				P1Games:          s.P1Games,
				P2Games:          s.P2Games,
				IsTiebreak:       s.IsTiebreak,
				P1TiebreakPoints: s.P1TiebreakPoints,
				P2TiebreakPoints: s.P2TiebreakPoints,
			}
			if err := repository.InsertSetScoreTx(tx, ss); err != nil {
				return fmt.Errorf("insert set score: %w", err)
			}
		}

		// 2. Update participant results.
		if err := repository.UpdateParticipantResultTx(tx, matchID, winnerPart.RegistrationID, "win"); err != nil {
			return err
		}
		if err := repository.UpdateParticipantResultTx(tx, matchID, loserPart.RegistrationID, "loss"); err != nil {
			return err
		}

		// 3. Complete the match.
		if err := repository.CompleteMatchTx(tx, matchID, winnerPart.RegistrationID, now); err != nil {
			return err
		}

		// 4. Update standings (group / round-robin stages only).
		if stage.Type == "group" || stage.Type == "round_robin" {
			if err := updateStandingsTx(tx, m, winnerPart.RegistrationID, loserPart.RegistrationID,
				winnerSets, loserSets, winnerGames, loserGames); err != nil {
				return err
			}
		}

		// 5. Advance winner to next match.
		if err := advanceWinner(tx, m, winnerPart.RegistrationID); err != nil {
			return err
		}

		// 6. Route loser to third-place match if configured.
		if m.LoserNextMatchID != nil && m.LoserNextMatchSlot != nil {
			lp := &model.MatchParticipant{
				MatchID:        *m.LoserNextMatchID,
				RegistrationID: loserPart.RegistrationID,
				Slot:           *m.LoserNextMatchSlot,
			}
			if err := repository.AddParticipantTx(tx, lp); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Return refreshed detail.
	sets, _ := e.matchRepo.GetSetScores(matchID)
	updatedParts, _ := e.matchRepo.GetParticipants(matchID)
	updatedMatch, _ := e.matchRepo.GetByID(matchID)
	return &model.MatchDetail{Match: updatedMatch, Participants: updatedParts, Sets: sets}, nil
}

// CorrectResult voids a completed result, reverses standings, and re-submits.
func (e *ResultEngine) CorrectResult(ctx context.Context, matchID string, req *model.SubmitResultRequest) (*model.MatchDetail, error) {
	m, err := e.matchRepo.GetByID(matchID)
	if err != nil {
		return nil, err
	}
	if m.Status != "completed" {
		return nil, fmt.Errorf("match is not completed, nothing to correct")
	}
	if m.NextMatchID != nil {
		next, err := e.matchRepo.GetByID(*m.NextMatchID)
		if err == nil && next.Status == "completed" {
			return nil, fmt.Errorf("next match already completed: correction blocked, requires admin resolution")
		}
	}

	stage, err := e.stageRepo.GetByID(m.StageID)
	if err != nil {
		return nil, err
	}

	err = e.store.RunInTx(ctx, func(tx *gorm.DB) error {
		// 1. Delete old set scores.
		if err := repository.DeleteSetScoresTx(tx, matchID); err != nil {
			return err
		}

		// 2. Reverse standings.
		if stage.Type == "group" || stage.Type == "round_robin" {
			parts, err := repository.GetParticipantsTx(tx, matchID)
			if err != nil {
				return err
			}
			if len(parts) == 2 {
				if err := reverseStandingsTx(tx, m, parts); err != nil {
					return err
				}
			}
		}

		// 3. Remove winner from next match slot (best-effort).
		if m.NextMatchID != nil && m.WinnerRegistrationID != nil {
			_ = repository.DeleteParticipantSlotTx(tx, *m.NextMatchID, *m.WinnerRegistrationID)
		}

		// 4. Reset match to pending.
		return repository.UpdateMatchStatusTx(tx, matchID, "pending")
	})
	if err != nil {
		return nil, err
	}

	return e.Submit(ctx, matchID, req)
}

// ── Standings helpers ──────────────────────────────────────────────────────

func updateStandingsTx(
	tx *gorm.DB,
	m *model.Match,
	winnerID, loserID string,
	winnerSets, loserSets, winnerGames, loserGames int,
) error {
	for _, upd := range []struct {
		id       string
		winner   bool
		sets     int
		oppSets  int
		games    int
		oppGames int
	}{
		{winnerID, true, winnerSets, loserSets, winnerGames, loserGames},
		{loserID, false, loserSets, winnerSets, loserGames, winnerGames},
	} {
		s, err := repository.GetStandingsByRegistrationTx(tx, m.StageID, m.GroupID, upd.id)
		if err != nil {
			s = &model.Standings{
				StageID:        m.StageID,
				GroupID:        m.GroupID,
				RegistrationID: upd.id,
			}
		}
		s.Played++
		s.SetsWon += upd.sets
		s.SetsLost += upd.oppSets
		s.GamesWon += upd.games
		s.GamesLost += upd.oppGames
		if upd.winner {
			s.Wins++
			s.Points += 2
		} else {
			s.Losses++
		}
		if err := repository.UpsertStandingsTx(tx, s); err != nil {
			return err
		}
	}
	return nil
}

func reverseStandingsTx(tx *gorm.DB, m *model.Match, participants []*model.MatchParticipant) error {
	for _, p := range participants {
		s, err := repository.GetStandingsByRegistrationTx(tx, m.StageID, m.GroupID, p.RegistrationID)
		if err != nil {
			continue
		}
		s.Played = max(0, s.Played-1)
		if p.Result != nil {
			switch *p.Result {
			case "win":
				s.Wins = max(0, s.Wins-1)
				s.Points = max(0, s.Points-2)
			case "loss":
				s.Losses = max(0, s.Losses-1)
			}
		}
		if err := repository.UpsertStandingsTx(tx, s); err != nil {
			return err
		}
	}
	return nil
}

// ── Validation ─────────────────────────────────────────────────────────────

func validateSets(sets []model.SetScoreInput, bestOf int) error {
	if len(sets) == 0 {
		return fmt.Errorf("at least one set is required")
	}
	if len(sets) > bestOf {
		return fmt.Errorf("too many sets: best of %d, got %d", bestOf, len(sets))
	}
	setsToWin := (bestOf / 2) + 1
	p1, p2 := 0, 0
	for i, s := range sets {
		if s.P1Games == s.P2Games {
			return fmt.Errorf("set %d: scores cannot be equal (%d-%d)", s.SetNumber, s.P1Games, s.P2Games)
		}
		if s.P1Games > s.P2Games {
			p1++
		} else {
			p2++
		}
		if i < len(sets)-1 && (p1 >= setsToWin || p2 >= setsToWin) {
			return fmt.Errorf("match already decided at set %d but more sets were submitted", i+1)
		}
	}
	return nil
}

func participantBySlot(parts []*model.MatchParticipant, slot int) *model.MatchParticipant {
	for _, p := range parts {
		if p.Slot == slot {
			return p
		}
	}
	return nil
}