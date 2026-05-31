package service

import (
	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByID(id int64) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) Update(user *model.User) error {
	return s.userRepo.Update(user)
}

// Address operations

func (s *UserService) ListAddresses(userID int64) ([]model.UserAddress, error) {
	return s.userRepo.ListAddresses(userID)
}

func (s *UserService) CreateAddress(addr *model.UserAddress) error {
	if addr.IsDefault {
		s.userRepo.ClearDefaultAddress(addr.UserID)
	}
	return s.userRepo.CreateAddress(addr)
}

func (s *UserService) UpdateAddress(addr *model.UserAddress) error {
	if addr.IsDefault {
		s.userRepo.ClearDefaultAddress(addr.UserID)
	}
	return s.userRepo.UpdateAddress(addr)
}

func (s *UserService) DeleteAddress(id, userID int64) error {
	return s.userRepo.DeleteAddress(id, userID)
}

func (s *UserService) GetAddress(id, userID int64) (*model.UserAddress, error) {
	return s.userRepo.GetAddress(id, userID)
}
