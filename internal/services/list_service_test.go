package services

import (
	"testing"
	"trello-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockListRepository struct {
	mock.Mock
}

func (m *MockListRepository) CreateList(list *models.List) error {
	args := m.Called(list)
	return args.Error(0)
}
func (m *MockListRepository) GetListsByBoardID(boardID uint) ([]models.List, error) {
	args := m.Called(boardID)
	return args.Get(0).([]models.List), args.Error(1)
}
func (m *MockListRepository) GetListByID(id uint) (*models.List, error) {
	args := m.Called(id)
	return args.Get(0).(*models.List), args.Error(1)
}
func (m *MockListRepository) UpdateList(list *models.List) error {
	args := m.Called(list)
	return args.Error(0)
}
func (m *MockListRepository) DeleteList(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestListService_CreateList(t *testing.T) {
	repo := new(MockListRepository)
	service := ListService{ListRepo: repo}
	boardID := uint(1)
	existing := []models.List{{Position: 0}, {Position: 1}}
	repo.On("GetListsByBoardID", boardID).Return(existing, nil)
	repo.On("CreateList", mock.MatchedBy(func(l *models.List) bool {
		return l.BoardID == boardID && l.Name == "NewList" && l.Position == len(existing)
	})).Return(nil)

	list, err := service.CreateList(boardID, "NewList")

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, boardID, list.BoardID)
	assert.Equal(t, "NewList", list.Name)
	assert.Equal(t, len(existing), list.Position)
}

func TestListService_GetLists(t *testing.T) {
	repo := new(MockListRepository)
	service := ListService{ListRepo: repo}
	boardID := uint(2)
	existing := []models.List{{ID: 1}, {ID: 2}}
	repo.On("GetListsByBoardID", boardID).Return(existing, nil)

	lists, err := service.GetLists(boardID)

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Len(t, lists, 2)
}

func TestListService_UpdateList(t *testing.T) {
	repo := new(MockListRepository)
	service := ListService{ListRepo: repo}
	id := uint(3)
	old := &models.List{ID: id, Name: "Old"}
	repo.On("GetListByID", id).Return(old, nil)
	repo.On("UpdateList", old).Return(nil)

	err := service.UpdateList(id, "New")

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, "New", old.Name)
}

func TestListService_DeleteList(t *testing.T) {
	repo := new(MockListRepository)
	service := ListService{ListRepo: repo}
	id := uint(4)
	repo.On("DeleteList", id).Return(nil)

	err := service.DeleteList(id)

	repo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestListService_MoveList_NoOp(t *testing.T) {
	repo := new(MockListRepository)
	service := ListService{ListRepo: repo}
	id := uint(5)
	list := &models.List{ID: id, Position: 2}
	repo.On("GetListByID", id).Return(list, nil)

	err := service.MoveList(id, 2)

	repo.AssertExpectations(t)
	assert.NoError(t, err)
}
