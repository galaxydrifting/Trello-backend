package repositories

import (
	"trello-backend/internal/models"

	"gorm.io/gorm"
)

type BoardRepository interface {
	CreateBoard(board *models.Board) error
	GetBoardByID(id uint) (*models.Board, error)
	UpdateBoard(board *models.Board) error
	DeleteBoard(id uint) error
	FindBoardsByUserID(userID string, boards *[]models.Board) error
}

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) CreateBoard(board *models.Board) error {
	return r.db.Create(board).Error
}

func (r *boardRepository) GetBoardByID(id uint) (*models.Board, error) {
	var board models.Board
	if err := r.db.Preload("Lists.Cards").First(&board, id).Error; err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) UpdateBoard(board *models.Board) error {
	return r.db.Save(board).Error
}

func (r *boardRepository) DeleteBoard(id uint) error {
	return r.db.Delete(&models.Board{}, id).Error
}

func (r *boardRepository) FindBoardsByUserID(userID string, boards *[]models.Board) error {
	return r.db.Preload("Lists.Cards").Where("user_id = ?", userID).Find(boards).Error
}
