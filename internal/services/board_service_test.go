package services

import (
	"testing"

	"trello-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBoardRepository struct {
	mock.Mock
}

func (m *MockBoardRepository) CreateBoard(board *models.Board) error {
	args := m.Called(board)
	return args.Error(0)
}

func (m *MockBoardRepository) GetBoardByID(id uint) (*models.Board, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Board), args.Error(1)
}

func (m *MockBoardRepository) UpdateBoard(board *models.Board) error {
	args := m.Called(board)
	return args.Error(0)
}

func (m *MockBoardRepository) DeleteBoard(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBoardService_CreateBoard(t *testing.T) {
	repo := new(MockBoardRepository)
	service := BoardService{BoardRepo: repo}

	board := &models.Board{Name: "Test Board"}
	repo.On("CreateBoard", board).Return(nil)

	result, err := service.CreateBoard("Test Board")

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, "Test Board", result.Name)
}

func TestBoardService_GetBoard(t *testing.T) {
	repo := new(MockBoardRepository)
	service := BoardService{BoardRepo: repo}

	board := &models.Board{ID: 1, Name: "Test Board"}
	repo.On("GetBoardByID", uint(1)).Return(board, nil)

	result, err := service.GetBoard(1)

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Test Board", result.Name)
}

func TestBoardService_UpdateBoard(t *testing.T) {
	repo := new(MockBoardRepository)
	service := BoardService{BoardRepo: repo}
	id := uint(10)
	old := &models.Board{ID: id, Name: "OldName"}
	repo.On("GetBoardByID", id).Return(old, nil)
	repo.On("UpdateBoard", old).Return(nil)

	err := service.UpdateBoard(id, "NewName")

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, "NewName", old.Name)
}

func TestBoardService_DeleteBoard(t *testing.T) {
	repo := new(MockBoardRepository)
	service := BoardService{BoardRepo: repo}
	id := uint(20)
	repo.On("DeleteBoard", id).Return(nil)

	err := service.DeleteBoard(id)

	repo.AssertExpectations(t)
	assert.NoError(t, err)
}
