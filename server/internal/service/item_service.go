package service

import (
	"github.com/neobarter/server/internal/model"
	"github.com/neobarter/server/internal/pkg/mq"
	"github.com/neobarter/server/internal/repository"
)

type ItemService struct {
	itemRepo  *repository.ItemRepository
	publisher *mq.Publisher // 可为 nil（MQ 未连接时降级）
}

func NewItemService(itemRepo *repository.ItemRepository, publisher *mq.Publisher) *ItemService {
	return &ItemService{itemRepo: itemRepo, publisher: publisher}
}

func (s *ItemService) Create(item *model.Item) error {
	if err := s.itemRepo.Create(item); err != nil {
		return err
	}
	s.publishEvent(mq.EventCreate, item)
	return nil
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
	if err := s.itemRepo.Update(item); err != nil {
		return err
	}
	s.publishEvent(mq.EventUpdate, item)
	return nil
}

func (s *ItemService) UpdateStatus(id int64, userID int64, status string) error {
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return ErrForbidden
	}
	if err := s.itemRepo.UpdateStatus(id, status); err != nil {
		return err
	}
	// 下架/删除时从搜索中移除
	if status == model.ItemStatusDeleted || status == model.ItemStatusInactive {
		s.publishDelete(id)
	} else {
		item.Status = status
		s.publishEvent(mq.EventUpdate, item)
	}
	return nil
}

func (s *ItemService) Delete(id int64, userID int64) error {
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		return err
	}
	if item.UserID != userID {
		return ErrForbidden
	}
	if err := s.itemRepo.UpdateStatus(id, model.ItemStatusDeleted); err != nil {
		return err
	}
	s.publishDelete(id)
	return nil
}

func (s *ItemService) List(q repository.ItemQuery) ([]model.Item, int64, error) {
	return s.itemRepo.List(q)
}

func (s *ItemService) ListCategories() ([]model.Category, error) {
	return s.itemRepo.ListCategories()
}

// publishEvent 发布物品索引事件到 MQ
func (s *ItemService) publishEvent(eventType mq.EventType, item *model.Item) {
	if s.publisher == nil {
		return
	}

	categoryName := ""
	if item.Category != nil {
		categoryName = item.Category.Name
	}
	userNickname := ""
	if item.User != nil {
		userNickname = item.User.Nickname
	}

	data := map[string]interface{}{
		"id":              item.ID,
		"user_id":         item.UserID,
		"title":           item.Title,
		"description":     item.Description,
		"category_id":     item.CategoryID,
		"category_name":   categoryName,
		"estimated_value": item.EstimatedValue,
		"condition":       item.Condition,
		"images":          item.Images,
		"status":          item.Status,
		"location":        item.Location,
		"view_count":      item.ViewCount,
		"want_items":      item.WantItems,
		"user_nickname":   userNickname,
		"created_at":      item.CreatedAt,
		"updated_at":      item.UpdatedAt,
	}

	s.publisher.PublishItemEvent(mq.ItemEvent{
		Type:   eventType,
		ItemID: item.ID,
		Data:   data,
	})
}

func (s *ItemService) publishDelete(itemID int64) {
	if s.publisher == nil {
		return
	}
	s.publisher.PublishItemEvent(mq.ItemEvent{
		Type:   mq.EventDelete,
		ItemID: itemID,
	})
}
