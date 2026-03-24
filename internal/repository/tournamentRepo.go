package repository

import (

	"github.com/uthmanduro/BracketForge/internal/model"
	"gorm.io/gorm"
)

type TournamentRepo struct {
	DB *gorm.DB
}

func NewTournamentRepo(db *gorm.DB) *TournamentRepo {
	return &TournamentRepo{
		DB: db,
	}
}

func (r *TournamentRepo) Create(t *model.Tournament) error {
	return r.DB.Create(t).Error
}
 
func (r *TournamentRepo) GetByID(id, orgID string) (*model.Tournament, error) {
	var t model.Tournament
	return &t, r.DB.First(&t, "id = ? AND organization_id = ?", id, orgID).Error
}
 
func (r *TournamentRepo) ListByOrg(orgID string) ([]*model.Tournament, error) {
	var list []*model.Tournament
	return list, r.DB.Where("organization_id = ?", orgID).
		Order("created_at DESC").Find(&list).Error
}
 
func (r *TournamentRepo) UpdateStatus(id, orgID, status string) error {
	return r.DB.Model(&model.Tournament{}).
		Where("id = ? AND organization_id = ?", id, orgID).
		Update("status", status).Error
}


// ── Registrations ──────────────────────────────────────────────────────────
 
func (r *TournamentRepo) Register(reg *model.TournamentRegistration) error {
	return r.DB.Create(reg).Error
}
 
func (r *TournamentRepo) GetRegistration(id string) (*model.TournamentRegistration, error) {
	var reg model.TournamentRegistration
	err := r.DB.Raw(`
		SELECT tr.*, p.name AS player_name
		FROM tournament_registrations tr
		JOIN players p ON p.id = tr.player_id
		WHERE tr.id = ?`, id).Scan(&reg).Error
	return &reg, err
}
 
func (r *TournamentRepo) GetRegistrationByPlayerAndTournament(playerID, tournamentID string) (*model.TournamentRegistration, error) {
	var reg model.TournamentRegistration
	return &reg, r.DB.First(&reg, "player_id = ? AND tournament_id = ?", playerID, tournamentID).Error
}
 
func (r *TournamentRepo) ListRegistrations(tournamentID string) ([]*model.TournamentRegistration, error) {
	var list []*model.TournamentRegistration
	err := r.DB.Raw(`
		SELECT tr.*, p.name AS player_name
		FROM tournament_registrations tr
		JOIN players p ON p.id = tr.player_id
		WHERE tr.tournament_id = ?
		ORDER BY tr.seed NULLS LAST, p.name`, tournamentID).Scan(&list).Error
	return list, err
}
 
func (r *TournamentRepo) ListActiveRegistrations(tournamentID string) ([]*model.TournamentRegistration, error) {
	var list []*model.TournamentRegistration
	err := r.DB.Raw(`
		SELECT tr.*, p.name AS player_name
		FROM tournament_registrations tr
		JOIN players p ON p.id = tr.player_id
		WHERE tr.tournament_id = ? AND tr.status IN ('registered','checked_in')
		ORDER BY tr.seed NULLS LAST, p.name`, tournamentID).Scan(&list).Error
	return list, err
}
 
func (r *TournamentRepo) UpdateRegistrationStatus(id, status string) error {
	return r.DB.Model(&model.TournamentRegistration{}).
		Where("id = ?", id).
		Update("status", status).Error
}
 