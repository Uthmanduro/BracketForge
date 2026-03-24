package database

import (
	"github.com/uthmanduro/BracketForge/internal/model"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	// AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes
	// It will change existing column's type if it's size, precision, nullable changed
	// AutoMigrate will not delete unused columns to protect your data
	models := []interface{}{
		&model.Organization{},
		&model.User{},
		&model.Player{},
		&model.Tournament{},
		&model.TournamentRegistration{},
		&model.Stage{},
		&model.Group{},
		&model.GroupRegistration{},
		&model.Match{},
		&model.MatchParticipant{},
		&model.SetScore{},
		&model.Standings{},

	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	return nil
}