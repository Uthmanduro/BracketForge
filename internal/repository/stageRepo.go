package repository

import (
	"github.com/uthmanduro/BracketForge/internal/model"

	"gorm.io/gorm"
)

type StageRepo struct{ db *gorm.DB }

func NewStageRepo(db *gorm.DB) *StageRepo { return &StageRepo{db} }

func (r *StageRepo) Create(s *model.Stage) error {
	return r.db.Create(s).Error
}

func (r *StageRepo) GetByID(id string) (*model.Stage, error) {
	var s model.Stage
	return &s, r.db.First(&s, "id = ?", id).Error
}

func (r *StageRepo) ListByTournament(tournamentID string) ([]*model.Stage, error) {
	var list []*model.Stage
	return list, r.db.Where("tournament_id = ?", tournamentID).
		Order("stage_order").Find(&list).Error
}

// ── Groups ─────────────────────────────────────────────────────────────────

func (r *StageRepo) CreateGroup(g *model.Group) error {
	return r.db.Create(g).Error
}

func (r *StageRepo) CreateGroupTx(tx *gorm.DB, g *model.Group) error {
	return tx.Create(g).Error
}

func (r *StageRepo) GetGroup(id string) (*model.Group, error) {
	var g model.Group
	return &g, r.db.First(&g, "id = ?", id).Error
}

func (r *StageRepo) ListGroups(stageID string) ([]*model.Group, error) {
	var list []*model.Group
	return list, r.db.Where("stage_id = ?", stageID).Order("name").Find(&list).Error
}

// ── Group registrations ────────────────────────────────────────────────────

func (r *StageRepo) AddGroupRegistration(gr *model.GroupRegistration) error {
	return r.db.Create(gr).Error
}

func (r *StageRepo) AddGroupRegistrationTx(tx *gorm.DB, gr *model.GroupRegistration) error {
	return tx.Create(gr).Error
}

func (r *StageRepo) ListGroupRegistrations(groupID string) ([]*model.GroupRegistration, error) {
	var list []*model.GroupRegistration
	return list, r.db.Where("group_id = ?", groupID).Find(&list).Error
}

func (r *StageRepo) WithTx(tx *gorm.DB) *StageRepo { return &StageRepo{db: tx} }