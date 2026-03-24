package repository

import (

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
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	return &user, r.db.Where("email = ?", email).First(&user).Error
}

func (r *UserRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) GetByID(id string) (*model.User, error) {
	var u model.User
	return &u, r.db.First(&u, "id = ?", id).Error
}