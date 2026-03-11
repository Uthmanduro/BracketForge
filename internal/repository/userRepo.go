package repository

import (
	"context"

	"github.com/uthmanduro/BracketForge/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(userID uint) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	db.AutoMigrate(&model.User{}) // Ensure the User table is created
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No user found with this email
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.Create(user).Error
}