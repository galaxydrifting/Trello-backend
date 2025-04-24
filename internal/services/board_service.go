package services

import (
	"trello-backend/internal/models"
	"trello-backend/internal/repositories"
)

type BoardService interface {
	CreateBoard(name string, userID string) (*models.Board, error)
	GetBoard(id uint) (*models.Board, error)
	UpdateBoard(id uint, name string) error
	DeleteBoard(id uint) error
	GetBoardsByUserID(userID string) ([]models.Board, error)
}

type boardService struct {
	boardRepo repositories.BoardRepository
}

func NewBoardService(repo repositories.BoardRepository) BoardService {
	return &boardService{boardRepo: repo}
}

func (s *boardService) CreateBoard(name string, userID string) (*models.Board, error) {
	board := &models.Board{Name: name, UserID: userID}
	if err := s.boardRepo.CreateBoard(board); err != nil {
		return nil, err
	}
	return board, nil
}

func (s *boardService) GetBoard(id uint) (*models.Board, error) {
	return s.boardRepo.GetBoardByID(id)
}

func (s *boardService) UpdateBoard(id uint, name string) error {
	board, err := s.boardRepo.GetBoardByID(id)
	if err != nil {
		return err
	}
	board.Name = name
	return s.boardRepo.UpdateBoard(board)
}

func (s *boardService) DeleteBoard(id uint) error {
	return s.boardRepo.DeleteBoard(id)
}

func (s *boardService) GetBoardsByUserID(userID string) ([]models.Board, error) {
	var boards []models.Board
	err := s.boardRepo.FindBoardsByUserID(userID, &boards)
	return boards, err
}
