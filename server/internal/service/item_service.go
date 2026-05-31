package service

import (
	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/repository"
)

type ItemService struct {
	itemRepo *repository.ItemRepository
}

func NewItemService(itemRepo *repository.ItemRepository) *ItemService {
	return &ItemService{itemRepo: itemRepo}
}

func (s *ItemService) Create(item *model.Item) error {
	return s.itemRepo.Create(item)
}

func (s *ItemService) GetByID(id int64) (*model.Item, error) {
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	// 增加浏览量
	s.itemRepo.IncrViewCount(id)
	return item, nil
}

func (s *ItemService) Update(item *model.Item) error {
	return s.itemRepo.Update(item)
}

func (s *ItemService) UpdateStatus(id int64, userID int64, status string) error {
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return ErrForbidden
	}
	return s.itemRepo.UpdateStatus(id, status)
}

func (s *ItemService) Delete(id int64, userID int64) error {
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return ErrForbidden
	}
	return s.itemRepo.UpdateStatus(id, model.ItemStatusDeleted)
}

func (s *ItemService) List(q repository.ItemQuery) ([]model.Item, int64, error) {
	return s.itemRepo.List(q)
}

func (s *ItemService) ListCategories() ([]model.Category, error) {
	return s.itemRepo.ListCategories()
}
