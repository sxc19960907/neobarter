package repository

import (
	"github.com/neobarter/server/internal/model"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(item *model.Item) error {
	return r.db.Create(item).Error
}

func (r *ItemRepository) GetByID(id int64) (*model.Item, error) {
	var item model.Item
	err := r.db.Preload("User").Preload("Category").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) Update(item *model.Item) error {
	return r.db.Save(item).Error
}

func (r *ItemRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&model.Item{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ItemRepository) IncrViewCount(id int64) error {
	return r.db.Model(&model.Item{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

type ItemQuery struct {
	UserID     int64
	CategoryID int
	Status     string
	Keyword    string
	Location   string
	Condition  string
	MinValue   float64
	MaxValue   float64
	SortBy     string // created_at / view_count / estimated_value
	Page       int
	PageSize   int
}

func (r *ItemRepository) List(q ItemQuery) ([]model.Item, int64, error) {
	var items []model.Item
	var total int64

	query := r.db.Model(&model.Item{})

	if q.UserID > 0 {
		query = query.Where("user_id = ?", q.UserID)
	}
	if q.CategoryID > 0 {
		query = query.Where("category_id = ?", q.CategoryID)
	}
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	} else {
		query = query.Where("status = ?", model.ItemStatusActive)
	}
	if q.Keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+q.Keyword+"%", "%"+q.Keyword+"%")
	}
	if q.Location != "" {
		query = query.Where("location ILIKE ?", "%"+q.Location+"%")
	}
	if q.Condition != "" {
		query = query.Where("condition = ?", q.Condition)
	}
	if q.MinValue > 0 {
		query = query.Where("estimated_value >= ?", q.MinValue)
	}
	if q.MaxValue > 0 {
		query = query.Where("estimated_value <= ?", q.MaxValue)
	}

	query.Count(&total)

	// 排序
	switch q.SortBy {
	case "view_count":
		query = query.Order("view_count DESC")
	case "estimated_value":
		query = query.Order("estimated_value DESC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Preload("User").Preload("Category").
		Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).
		Find(&items).Error

	return items, total, err
}

func (r *ItemRepository) ListCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}
