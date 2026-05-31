package repository

import (
	"github.com/neobarter/server/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) UpdateCreditScore(userID int64, delta int) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("credit_score", gorm.Expr("credit_score + ?", delta)).Error
}

// Address operations

func (r *UserRepository) CreateAddress(addr *model.UserAddress) error {
	return r.db.Create(addr).Error
}

func (r *UserRepository) ListAddresses(userID int64) ([]model.UserAddress, error) {
	var addresses []model.UserAddress
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, id DESC").Find(&addresses).Error
	return addresses, err
}

func (r *UserRepository) GetAddress(id, userID int64) (*model.UserAddress, error) {
	var addr model.UserAddress
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&addr).Error
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func (r *UserRepository) UpdateAddress(addr *model.UserAddress) error {
	return r.db.Save(addr).Error
}

func (r *UserRepository) DeleteAddress(id, userID int64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.UserAddress{}).Error
}

func (r *UserRepository) ClearDefaultAddress(userID int64) error {
	return r.db.Model(&model.UserAddress{}).Where("user_id = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error
}
