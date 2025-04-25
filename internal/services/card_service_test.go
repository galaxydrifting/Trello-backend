package services

import (
	"errors"
	"testing"
	"trello-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCardRepository struct {
	mock.Mock
}

func (m *MockCardRepository) CreateCard(card *models.Card) error {
	args := m.Called(card)
	return args.Error(0)
}
func (m *MockCardRepository) GetCardsByListID(listID uint) ([]models.Card, error) {
	args := m.Called(listID)
	return args.Get(0).([]models.Card), args.Error(1)
}
func (m *MockCardRepository) GetCardByID(id uint) (*models.Card, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Card), args.Error(1)
}
func (m *MockCardRepository) UpdateCard(card *models.Card) error {
	args := m.Called(card)
	return args.Error(0)
}
func (m *MockCardRepository) DeleteCard(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCardService_CreateCard(t *testing.T) {
	repo := new(MockCardRepository)
	service := NewCardService(repo)
	listID := uint(1)
	existing := []models.Card{{Position: 0}, {Position: 1}}
	repo.On("GetCardsByListID", listID).Return(existing, nil)
	repo.On("CreateCard", mock.MatchedBy(func(c *models.Card) bool {
		return c.ListID == listID && c.Title == "Title" && c.Content == "Content" && c.Position == len(existing)
	})).Return(nil)

	card, err := service.CreateCard(listID, "Title", "Content")

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, listID, card.ListID)
}

func TestCardService_GetCards(t *testing.T) {
	repo := new(MockCardRepository)
	service := NewCardService(repo)
	listID := uint(2)
	existing := []models.Card{{ID: 1}, {ID: 2}}
	repo.On("GetCardsByListID", listID).Return(existing, nil)

	cards, err := service.GetCards(listID)

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Len(t, cards, 2)
}

func TestCardService_UpdateCard(t *testing.T) {
	repo := new(MockCardRepository)
	service := NewCardService(repo)
	id := uint(3)
	old := &models.Card{ID: id, Title: "Old"}
	repo.On("GetCardByID", id).Return(old, nil)
	repo.On("UpdateCard", old).Return(nil)

	err := service.UpdateCard(id, "New", "")

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, "New", old.Title)
}

func TestCardService_DeleteCard(t *testing.T) {
	repo := new(MockCardRepository)
	service := NewCardService(repo)
	id := uint(4)
	repo.On("DeleteCard", id).Return(nil)

	err := service.DeleteCard(id)

	repo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestCardService_MoveCard_GetError(t *testing.T) {
	repo := new(MockCardRepository)
	service := NewCardService(repo)
	id := uint(5)
	repo.On("GetCardByID", id).Return(nil, errors.New("not found"))

	err := service.MoveCard(id, 1, 0)

	repo.AssertExpectations(t)
	assert.Error(t, err)
}
