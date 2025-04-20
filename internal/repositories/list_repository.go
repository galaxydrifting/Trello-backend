package repositories

import (
	"trello-backend/internal/models"

	"gorm.io/gorm"
)

type ListRepository interface {
	CreateList(list *models.List) error
	GetListsByBoardID(boardID uint) ([]models.List, error)
	GetListByID(id uint) (*models.List, error)
	UpdateList(list *models.List) error
	DeleteList(id uint) error
}

type listRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) ListRepository {
	return &listRepository{db: db}
}

func (r *listRepository) CreateList(list *models.List) error {
	return r.db.Create(list).Error
}

func (r *listRepository) GetListsByBoardID(boardID uint) ([]models.List, error) {
	var lists []models.List
	err := r.db.Where("board_id = ?", boardID).Order("position").Find(&lists).Error
	return lists, err
}

func (r *listRepository) GetListByID(id uint) (*models.List, error) {
	var list models.List
	if err := r.db.First(&list, id).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *listRepository) UpdateList(list *models.List) error {
	return r.db.Save(list).Error
}

func (r *listRepository) DeleteList(id uint) error {
	return r.db.Delete(&models.List{}, id).Error
}
