package engine

import (
	"fmt"
	"time"

	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"

	"gorm.io/gorm"
)

// ── BYE engine ─────────────────────────────────────────────────────────────

type ByeEngine struct{ matchRepo *repository.MatchRepo }

func NewByeEngine(mr *repository.MatchRepo) *ByeEngine { return &ByeEngine{mr} }

func (e *ByeEngine) AutoAdvance(m *model.Match, participant *model.MatchParticipant) error {
	now := time.Now()
	if err := e.matchRepo.UpdateParticipantResult(m.ID, participant.RegistrationID, "bye_advance"); err != nil {
		return err
	}
	if err := e.matchRepo.CompleteMatch(m.ID, participant.RegistrationID, now); err != nil {
		return err
	}
	return advanceWinner(e.matchRepo.DB(), m, participant.RegistrationID)
}

// ── Walkover engine ────────────────────────────────────────────────────────

type WalkoverEngine struct{ matchRepo *repository.MatchRepo }

func NewWalkoverEngine(mr *repository.MatchRepo) *WalkoverEngine { return &WalkoverEngine{mr} }

func (e *WalkoverEngine) ProcessWithdrawal(registrationID string) error {
	pending, err := e.matchRepo.GetParticipantsByRegistration(registrationID)
	if err != nil {
		return err
	}
	for _, part := range pending {
		m, err := e.matchRepo.GetByID(part.MatchID)
		if err != nil {
			return err
		}
		participants, err := e.matchRepo.GetParticipants(m.ID)
		if err != nil {
			return err
		}

		var opponentID string
		for _, p := range participants {
			if p.RegistrationID != registrationID {
				opponentID = p.RegistrationID
				break
			}
		}

		if opponentID == "" {
			if err := e.matchRepo.UpdateStatus(m.ID, "cancelled"); err != nil {
				return err
			}
			continue
		}

		if err := e.matchRepo.UpdateParticipantResult(m.ID, registrationID, "walkover_loss"); err != nil {
			return err
		}
		if err := e.matchRepo.UpdateParticipantResult(m.ID, opponentID, "walkover_win"); err != nil {
			return err
		}
		if err := e.matchRepo.WalkoverMatch(m.ID, opponentID); err != nil {
			return err
		}
		if err := advanceWinner(e.matchRepo.DB(), m, opponentID); err != nil {
			return err
		}
	}
	return nil
}

// ── Bracket advancement (shared) ──────────────────────────────────────────

// advanceWinner places the winner into the next match slot.
// Works both inside and outside a transaction — caller passes the appropriate *gorm.DB.
func advanceWinner(db *gorm.DB, m *model.Match, winnerRegID string) error {
	if m.NextMatchID == nil {
		return nil // terminal match
	}
	if m.NextMatchSlot == nil {
		return fmt.Errorf("match %s: next_match_id set but next_match_slot is nil", m.ID)
	}

	p := &model.MatchParticipant{
		MatchID:        *m.NextMatchID,
		RegistrationID: winnerRegID,
		Slot:           *m.NextMatchSlot,
	}
	if err := repository.AddParticipantTx(db, p); err != nil {
		return err
	}

	// Mark next match pending once both slots are filled.
	parts, err := repository.GetParticipantsTx(db, *m.NextMatchID)
	if err != nil {
		return err
	}
	if len(parts) == 2 {
		return repository.UpdateMatchStatusTx(db, *m.NextMatchID, "pending")
	}
	return nil
}