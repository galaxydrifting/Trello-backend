package services

import (
	"trello-backend/internal/models"
)

type ListRepository interface {
	CreateList(list *models.List) error
	GetListsByBoardID(boardID uint) ([]models.List, error)
	GetListByID(id uint) (*models.List, error)
	UpdateList(list *models.List) error
	DeleteList(id uint) error
}

type ListService struct {
	ListRepo ListRepository
}

func (s *ListService) CreateList(boardID uint, name string) (*models.List, error) {
	lists, err := s.ListRepo.GetListsByBoardID(boardID)
	if err != nil {
		return nil, err
	}
	position := len(lists)
	list := &models.List{BoardID: boardID, Name: name, Position: position}
	if err := s.ListRepo.CreateList(list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *ListService) GetLists(boardID uint) ([]models.List, error) {
	return s.ListRepo.GetListsByBoardID(boardID)
}

func (s *ListService) UpdateList(id uint, name string) error {
	list, err := s.ListRepo.GetListByID(id)
	if err != nil {
		return err
	}
	list.Name = name
	return s.ListRepo.UpdateList(list)
}

func (s *ListService) DeleteList(id uint) error {
	return s.ListRepo.DeleteList(id)
}

func (s *ListService) MoveList(id uint, newPosition int) error {
	list, err := s.ListRepo.GetListByID(id)
	if err != nil {
		return err
	}
	oldPos := list.Position
	if newPosition == oldPos {
		return nil
	}
	lists, err := s.ListRepo.GetListsByBoardID(list.BoardID)
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
				if err := s.ListRepo.UpdateList(&l); err != nil {
					return err
				}
			}
		} else {
			if l.Position >= newPosition && l.Position < oldPos {
				l.Position++
				if err := s.ListRepo.UpdateList(&l); err != nil {
					return err
				}
			}
		}
	}
	list.Position = newPosition
	return s.ListRepo.UpdateList(list)
}
