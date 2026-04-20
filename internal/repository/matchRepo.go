package repository

import (
	"time"

	"github.com/uthmanduro/BracketForge/internal/model"

	"gorm.io/gorm"
)

type MatchRepo struct{ db *gorm.DB }

func NewMatchRepo(db *gorm.DB) *MatchRepo { return &MatchRepo{db} }

// ── Match CRUD ─────────────────────────────────────────────────────────────

func (r *MatchRepo) Create(m *model.Match) error {
	return r.db.Create(m).Error
}

func (r *MatchRepo) GetByID(id string) (*model.Match, error) {
	var m model.Match
	return &m, r.db.First(&m, "id = ?", id).Error
}

func (r *MatchRepo) ListByStage(stageID string) ([]*model.Match, error) {
	var list []*model.Match
	return list, r.db.Where("stage_id = ?", stageID).
		Order("round, match_position").Find(&list).Error
}

func (r *MatchRepo) ListByGroup(groupID string) ([]*model.Match, error) {
	var list []*model.Match
	return list, r.db.Where("group_id = ?", groupID).
		Order("round, match_position").Find(&list).Error
}

func (r *MatchRepo) UpdateStatus(id, status string) error {
	return r.db.Model(&model.Match{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *MatchRepo) CompleteMatch(id, winnerRegID string, completedAt time.Time) error {
	return r.db.Model(&model.Match{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":                 "completed",
		"winner_registration_id": winnerRegID,
		"completed_at":           completedAt,
	}).Error
}

func (r *MatchRepo) WalkoverMatch(id, winnerRegID string) error {
	return r.db.Model(&model.Match{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":                 "walkover",
		"winner_registration_id": winnerRegID,
		"completed_at":           time.Now(),
	}).Error
}

func (r *MatchRepo) ScheduleMatch(id string, scheduledAt time.Time) error {
	return r.db.Model(&model.Match{}).Where("id = ?", id).
		Update("scheduled_at", scheduledAt).Error
}

// ── Participants ───────────────────────────────────────────────────────────

func (r *MatchRepo) AddParticipant(p *model.MatchParticipant) error {
	return r.db.Create(p).Error
}

func (r *MatchRepo) GetParticipants(matchID string) ([]*model.MatchParticipant, error) {
	var list []*model.MatchParticipant
	err := r.db.Raw(`
		SELECT mp.*, p.name AS player_name
		FROM match_participants mp
		JOIN tournament_registrations tr ON tr.id = mp.registration_id
		JOIN players p ON p.id = tr.player_id
		WHERE mp.match_id = ?
		ORDER BY mp.slot`, matchID).Scan(&list).Error
	return list, err
}

func (r *MatchRepo) UpdateParticipantResult(matchID, registrationID, result string) error {
	return r.db.Model(&model.MatchParticipant{}).
		Where("match_id = ? AND registration_id = ?", matchID, registrationID).
		Update("result", result).Error
}

func (r *MatchRepo) GetParticipantsByRegistration(registrationID string) ([]*model.MatchParticipant, error) {
	var list []*model.MatchParticipant
	err := r.db.Raw(`
		SELECT mp.* FROM match_participants mp
		JOIN matches m ON m.id = mp.match_id
		WHERE mp.registration_id = ? AND m.status = 'pending'`, registrationID).
		Scan(&list).Error
	return list, err
}

func (r *MatchRepo) DeleteParticipantSlot(matchID, registrationID string) error {
	return r.db.Delete(&model.MatchParticipant{},
		"match_id = ? AND registration_id = ?", matchID, registrationID).Error
}

// ── Set scores ─────────────────────────────────────────────────────────────

func (r *MatchRepo) InsertSetScore(s *model.SetScore) error {
	return r.db.Create(s).Error
}

func (r *MatchRepo) GetSetScores(matchID string) ([]*model.SetScore, error) {
	var list []*model.SetScore
	return list, r.db.Where("match_id = ?", matchID).
		Order("set_number").Find(&list).Error
}

func (r *MatchRepo) DeleteSetScores(matchID string) error {
	return r.db.Delete(&model.SetScore{}, "match_id = ?", matchID).Error
}

// ── Mutual match helpers (tiebreaker) ──────────────────────────────────────

// GetMutualMatches returns completed matches where ALL participants are in registrationIDs.
func (r *MatchRepo) GetMutualMatches(groupID string, registrationIDs []string) ([]*model.Match, error) {
	var list []*model.Match
	err := r.db.Raw(`
		SELECT m.* FROM matches m
		WHERE m.group_id = ?
		  AND m.status = 'completed'
		  AND (
		    SELECT COUNT(*) FROM match_participants mp
		    WHERE mp.match_id = m.id AND mp.registration_id IN ?
		  ) = 2`,
		groupID, registrationIDs).Scan(&list).Error
	return list, err
}

// GetMutualSetStats aggregates set/game tallies from mutual matches only.
// Used for tiebreaker steps 3 and 4.
func (r *MatchRepo) GetMutualSetStats(groupID string, registrationIDs []string) (map[string]*model.MutualSetStats, error) {
	matches, err := r.GetMutualMatches(groupID, registrationIDs)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]*model.MutualSetStats, len(registrationIDs))
	for _, id := range registrationIDs {
		stats[id] = &model.MutualSetStats{}
	}

	for _, m := range matches {
		parts, err := r.GetParticipants(m.ID)
		if err != nil {
			return nil, err
		}
		sets, err := r.GetSetScores(m.ID)
		if err != nil {
			return nil, err
		}

		p1Sets, p2Sets, p1Games, p2Games := 0, 0, 0, 0
		for _, s := range sets {
			if s.P1Games > s.P2Games {
				p1Sets++
			} else {
				p2Sets++
			}
			p1Games += s.P1Games
			p2Games += s.P2Games
		}

		for _, p := range parts {
			st, ok := stats[p.RegistrationID]
			if !ok {
				continue
			}
			if p.Slot == 1 {
				st.SetsWon += p1Sets
				st.SetsLost += p2Sets
				st.GamesWon += p1Games
				st.GamesLost += p2Games
			} else {
				st.SetsWon += p2Sets
				st.SetsLost += p1Sets
				st.GamesWon += p2Games
				st.GamesLost += p1Games
			}
		}
	}
	return stats, nil
}

// ── Tx helpers (used inside store.RunInTx) ─────────────────────────────────

func (r *MatchRepo) WithTx(tx *gorm.DB) *MatchRepo { return &MatchRepo{db: tx} }

func CreateMatchTx(tx *gorm.DB, m *model.Match) error {
	return tx.Create(m).Error
}

func AddParticipantTx(tx *gorm.DB, p *model.MatchParticipant) error {
	return tx.Create(p).Error
}

func GetParticipantsTx(tx *gorm.DB, matchID string) ([]*model.MatchParticipant, error) {
	var list []*model.MatchParticipant
	err := tx.Raw(`
		SELECT mp.*, p.name AS player_name
		FROM match_participants mp
		JOIN tournament_registrations tr ON tr.id = mp.registration_id
		JOIN players p ON p.id = tr.player_id
		WHERE mp.match_id = ?
		ORDER BY mp.slot`, matchID).Scan(&list).Error
	return list, err
}

func UpdateParticipantResultTx(tx *gorm.DB, matchID, registrationID, result string) error {
	return tx.Model(&model.MatchParticipant{}).
		Where("match_id = ? AND registration_id = ?", matchID, registrationID).
		Update("result", result).Error
}

func CompleteMatchTx(tx *gorm.DB, id, winnerRegID string, completedAt time.Time) error {
	return tx.Model(&model.Match{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":                 "completed",
		"winner_registration_id": winnerRegID,
		"completed_at":           completedAt,
	}).Error
}

func UpdateMatchStatusTx(tx *gorm.DB, id, status string) error {
	return tx.Model(&model.Match{}).Where("id = ?", id).Update("status", status).Error
}

func InsertSetScoreTx(tx *gorm.DB, s *model.SetScore) error {
	return tx.Create(s).Error
}

func DeleteSetScoresTx(tx *gorm.DB, matchID string) error {
	return tx.Delete(&model.SetScore{}, "match_id = ?", matchID).Error
}

func DeleteParticipantSlotTx(tx *gorm.DB, matchID, registrationID string) error {
	return tx.Delete(&model.MatchParticipant{},
		"match_id = ? AND registration_id = ?", matchID, registrationID).Error
}

// DB returns the underlying *gorm.DB — used by engine helpers that need to
// pass the same connection (or transaction) to advanceWinner.
func (r *MatchRepo) DB() *gorm.DB { return r.db }