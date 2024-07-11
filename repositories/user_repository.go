package repositories

import (
	"desafio-picpay/models"
	"github.com/jinzhu/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}

func (r *GormUserRepository) FindByCPF(cpf string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("cpf = ?", cpf).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetInstance() *gorm.DB {
	return r.db
}
