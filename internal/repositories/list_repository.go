package repositories

import (
	"trello-backend/internal/models"

	"gorm.io/gorm"
)

// ListRepository 處理清單的資料操作
type ListRepository struct {
	DB *gorm.DB
}

func (r *ListRepository) CreateList(list *models.List) error {
	return r.DB.Create(list).Error
}

func (r *ListRepository) GetListsByBoardID(boardID uint) ([]models.List, error) {
	var lists []models.List
	err := r.DB.Where("board_id = ?", boardID).Order("position").Find(&lists).Error
	return lists, err
}

func (r *ListRepository) GetListByID(id uint) (*models.List, error) {
	var list models.List
	if err := r.DB.First(&list, id).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *ListRepository) UpdateList(list *models.List) error {
	return r.DB.Save(list).Error
}

func (r *ListRepository) DeleteList(id uint) error {
	return r.DB.Delete(&models.List{}, id).Error
}
