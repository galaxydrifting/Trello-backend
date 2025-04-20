package services

import (
	"trello-backend/internal/models"
	"trello-backend/internal/repositories"
)

type ListService interface {
	CreateList(boardID uint, name string) (*models.List, error)
	GetLists(boardID uint) ([]models.List, error)
	GetListByID(id uint) (*models.List, error)
	UpdateList(id uint, name string) error
	DeleteList(id uint) error
	MoveList(id uint, newPosition int) error
}

type listService struct {
	listRepo repositories.ListRepository
}

func NewListService(repo repositories.ListRepository) ListService {
	return &listService{listRepo: repo}
}

func (s *listService) CreateList(boardID uint, name string) (*models.List, error) {
	lists, err := s.listRepo.GetListsByBoardID(boardID)
	if err != nil {
		return nil, err
	}
	position := len(lists)
	list := &models.List{BoardID: boardID, Name: name, Position: position}
	if err := s.listRepo.CreateList(list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *listService) GetLists(boardID uint) ([]models.List, error) {
	return s.listRepo.GetListsByBoardID(boardID)
}

func (s *listService) GetListByID(id uint) (*models.List, error) {
	return s.listRepo.GetListByID(id)
}

func (s *listService) UpdateList(id uint, name string) error {
	list, err := s.listRepo.GetListByID(id)
	if err != nil {
		return err
	}
	list.Name = name
	return s.listRepo.UpdateList(list)
}

func (s *listService) DeleteList(id uint) error {
	return s.listRepo.DeleteList(id)
}

func (s *listService) MoveList(id uint, newPosition int) error {
	list, err := s.listRepo.GetListByID(id)
	if err != nil {
		return err
	}
	oldPos := list.Position
	if newPosition == oldPos {
		return nil
	}
	lists, err := s.listRepo.GetListsByBoardID(list.BoardID)
	if err != nil {
		return err
	}
	for _, l := range lists {
		if l.ID == id {
			continue
		}
		if oldPos < newPosition {
			if l.Position > oldPos && l.Position <= newPosition {
				l.Position--
				if err := s.listRepo.UpdateList(&l); err != nil {
					return err
				}
			}
		} else {
			if l.Position >= newPosition && l.Position < oldPos {
				l.Position++
				if err := s.listRepo.UpdateList(&l); err != nil {
					return err
				}
			}
		}
	}
	list.Position = newPosition
	return s.listRepo.UpdateList(list)
}
