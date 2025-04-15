package repositories

import (
	"gorm.io/gorm"
	"trello-backend/internal/models"
)

// BoardRepository handles CRUD operations for Board
type BoardRepository struct {
	DB *gorm.DB
}

func (r *BoardRepository) CreateBoard(board *models.Board) error {
	return r.DB.Create(board).Error
}

func (r *BoardRepository) GetBoardByID(id uint) (*models.Board, error) {
	var board models.Board
	if err := r.DB.Preload("Lists.Cards").First(&board, id).Error; err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *BoardRepository) UpdateBoard(board *models.Board) error {
	return r.DB.Save(board).Error
}

func (r *BoardRepository) DeleteBoard(id uint) error {
	return r.DB.Delete(&models.Board{}, id).Error
}

// Similar repositories can be created for List and Card.