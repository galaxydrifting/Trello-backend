package repositories

import (
	"trello-backend/internal/models"

	"gorm.io/gorm"
)

// CardRepository 處理卡片的資料操作
type CardRepository struct {
	DB *gorm.DB
}

func (r *CardRepository) CreateCard(card *models.Card) error {
	return r.DB.Create(card).Error
}

func (r *CardRepository) GetCardsByListID(listID uint) ([]models.Card, error) {
	var cards []models.Card
	err := r.DB.Where("list_id = ?", listID).Order("position").Find(&cards).Error
	return cards, err
}

func (r *CardRepository) GetCardByID(id uint) (*models.Card, error) {
	var card models.Card
	if err := r.DB.First(&card, id).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepository) UpdateCard(card *models.Card) error {
	return r.DB.Save(card).Error
}

func (r *CardRepository) DeleteCard(id uint) error {
	return r.DB.Delete(&models.Card{}, id).Error
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{DB: db}
}
