package repositories

import (
	"trello-backend/internal/models"

	"gorm.io/gorm"
)

type CardRepository interface {
	CreateCard(card *models.Card) error
	GetCardsByListID(listID uint) ([]models.Card, error)
	GetCardByID(id uint) (*models.Card, error)
	UpdateCard(card *models.Card) error
	DeleteCard(id uint) error
	GetCardsByBoardID(boardID uint) ([]models.Card, error)
}

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) CreateCard(card *models.Card) error {
	return r.db.Create(card).Error
}

func (r *cardRepository) GetCardsByListID(listID uint) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.Where("list_id = ?", listID).Order("position").Find(&cards).Error
	return cards, err
}

func (r *cardRepository) GetCardByID(id uint) (*models.Card, error) {
	var card models.Card
	if err := r.db.First(&card, id).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *cardRepository) UpdateCard(card *models.Card) error {
	return r.db.Save(card).Error
}

func (r *cardRepository) DeleteCard(id uint) error {
	return r.db.Delete(&models.Card{}, id).Error
}

func (r *cardRepository) GetCardsByBoardID(boardID uint) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.Where("board_id = ?", boardID).Order("position").Find(&cards).Error
	return cards, err
}
