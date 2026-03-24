package repository

import (
	"time"

	"github.com/uthmanduro/BracketForge/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StandingsRepo struct{ db *gorm.DB }

func NewStandingsRepo(db *gorm.DB) *StandingsRepo { return &StandingsRepo{db} }

// Upsert writes a standings row using GORM's OnConflict clause.
// Because group_id is nullable we use two separate upsert paths
// to match the two partial unique indexes defined in migrations.
func (r *StandingsRepo) Upsert(s *model.Standings) error {
	return upsertStandings(r.db, s)
}

func (r *StandingsRepo) GetByRegistration(stageID string, groupID *string, registrationID string) (*model.Standings, error) {
	var s model.Standings
	q := r.db.Where("stage_id = ? AND registration_id = ?", stageID, registrationID)
	if groupID == nil {
		q = q.Where("group_id IS NULL")
	} else {
		q = q.Where("group_id = ?", *groupID)
	}
	return &s, q.First(&s).Error
}

func (r *StandingsRepo) InitForRegistration(stageID string, groupID *string, registrationID string) error {
	s := &model.Standings{
		StageID:        stageID,
		GroupID:        groupID,
		RegistrationID: registrationID,
	}
	return upsertStandings(r.db, s)
}

func (r *StandingsRepo) ListByStage(stageID string) ([]*model.Standings, error) {
	var list []*model.Standings
	err := r.db.Raw(`
		SELECT s.*, p.name AS player_name
		FROM standings s
		JOIN tournament_registrations tr ON tr.id = s.registration_id
		JOIN players p ON p.id = tr.player_id
		WHERE s.stage_id = ?
		ORDER BY s.group_id NULLS FIRST, s.points DESC, s.sets_won DESC`, stageID).
		Scan(&list).Error
	return list, err
}

func (r *StandingsRepo) ListByGroup(groupID string) ([]*model.Standings, error) {
	var list []*model.Standings
	err := r.db.Raw(`
		SELECT s.*, p.name AS player_name
		FROM standings s
		JOIN tournament_registrations tr ON tr.id = s.registration_id
		JOIN players p ON p.id = tr.player_id
		WHERE s.group_id = ?
		ORDER BY s.points DESC, s.sets_won DESC`, groupID).
		Scan(&list).Error
	return list, err
}

// ── Tx helpers ─────────────────────────────────────────────────────────────

func UpsertStandingsTx(tx *gorm.DB, s *model.Standings) error {
	return upsertStandings(tx, s)
}

func GetStandingsByRegistrationTx(tx *gorm.DB, stageID string, groupID *string, registrationID string) (*model.Standings, error) {
	var s model.Standings
	q := tx.Where("stage_id = ? AND registration_id = ?", stageID, registrationID)
	if groupID == nil {
		q = q.Where("group_id IS NULL")
	} else {
		q = q.Where("group_id = ?", *groupID)
	}
	return &s, q.First(&s).Error
}

// ── shared upsert logic ────────────────────────────────────────────────────

func upsertStandings(db *gorm.DB, s *model.Standings) error {
	s.UpdatedAt = time.Now()
	// GORM OnConflict with DoUpdates handles INSERT ... ON CONFLICT DO UPDATE.
	// We specify the constraint columns explicitly since group_id is nullable.
	return db.Clauses(clause.OnConflict{
		// Postgres will match on the appropriate partial unique index.
		Columns: []clause.Column{
			{Name: "stage_id"},
			{Name: "registration_id"},
			{Name: "group_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"played", "wins", "losses", "points",
			"sets_won", "sets_lost", "games_won", "games_lost", "updated_at",
		}),
	}).Create(s).Error
}