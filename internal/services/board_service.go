package services

import (
	"trello-backend/internal/models"
)

// BoardRepository defines the interface for board data operations
type BoardRepository interface {
	CreateBoard(board *models.Board) error
	GetBoardByID(id uint) (*models.Board, error)
	UpdateBoard(board *models.Board) error
	DeleteBoard(id uint) error
}

// BoardService provides business logic for boards
type BoardService struct {
	BoardRepo BoardRepository
}

func (s *BoardService) CreateBoard(name string) (*models.Board, error) {
	board := &models.Board{Name: name}
	if err := s.BoardRepo.CreateBoard(board); err != nil {
		return nil, err
	}
	return board, nil
}

func (s *BoardService) GetBoard(id uint) (*models.Board, error) {
	return s.BoardRepo.GetBoardByID(id)
}

func (s *BoardService) UpdateBoard(id uint, name string) error {
	board, err := s.BoardRepo.GetBoardByID(id)
	if err != nil {
		return err
	}
	board.Name = name
	return s.BoardRepo.UpdateBoard(board)
}

func (s *BoardService) DeleteBoard(id uint) error {
	return s.BoardRepo.DeleteBoard(id)
}