package services

import (
	"trello-backend/internal/models"
)

type CardRepository interface {
	CreateCard(card *models.Card) error
	GetCardsByListID(listID uint) ([]models.Card, error)
	GetCardByID(id uint) (*models.Card, error)
	UpdateCard(card *models.Card) error
	DeleteCard(id uint) error
}

type CardService struct {
	CardRepo CardRepository
}

func (s *CardService) CreateCard(listID uint, title, content string) (*models.Card, error) {
	cards, err := s.CardRepo.GetCardsByListID(listID)
	if err != nil {
		return nil, err
	}
	position := len(cards)
	card := &models.Card{ListID: listID, Title: title, Content: content, Position: position}
	if err := s.CardRepo.CreateCard(card); err != nil {
		return nil, err
	}
	return card, nil
}

func (s *CardService) GetCards(listID uint) ([]models.Card, error) {
	return s.CardRepo.GetCardsByListID(listID)
}

func (s *CardService) UpdateCard(id uint, title, content string) error {
	card, err := s.CardRepo.GetCardByID(id)
	if err != nil {
		return err
	}
	card.Title = title
	card.Content = content
	return s.CardRepo.UpdateCard(card)
}

func (s *CardService) DeleteCard(id uint) error {
	return s.CardRepo.DeleteCard(id)
}

func (s *CardService) MoveCard(id, targetListID uint, newPosition int) error {
	card, err := s.CardRepo.GetCardByID(id)
	if err != nil {
		return err
	}
	oldListID := card.ListID
	oldPos := card.Position
	// 調整原清單中的卡片位置
	oldCards, err := s.CardRepo.GetCardsByListID(oldListID)
	if err != nil {
		return err
	}
	for _, c := range oldCards {
		if c.ID == id {
			continue
		}
		if c.Position > oldPos {
			c.Position--
			if err := s.CardRepo.UpdateCard(&c); err != nil {
				return err
			}
		}
	}
	// 插入到新清單
	card.ListID = targetListID
	newCards, err := s.CardRepo.GetCardsByListID(targetListID)
	if err != nil {
		return err
	}
	if newPosition > len(newCards) {
		newPosition = len(newCards)
	}
	for _, c := range newCards {
		if c.Position >= newPosition {
			c.Position++
			if err := s.CardRepo.UpdateCard(&c); err != nil {
				return err
			}
		}
	}
	card.Position = newPosition
	return s.CardRepo.UpdateCard(card)
}
