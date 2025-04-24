package graph

import (
	"context"
	"errors"
	"strconv"
	"trello-backend/graph/model"
	"trello-backend/pkg/utils"
)

// Board 相關 resolver function

func (r *mutationResolver) CreateBoard(ctx context.Context, input model.CreateBoardInput) (*model.Board, error) {
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		return nil, errors.New("未驗證身份")
	}
	b, err := r.BoardService.CreateBoard(input.Name, userID)
	if err != nil {
		return nil, err
	}
	return &model.Board{
		ID:        strconv.FormatUint(uint64(b.ID), 10),
		Name:      b.Name,
		CreatedAt: b.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: b.UpdatedAt.Format(utils.TimeFormat),
	}, nil
}

func (r *mutationResolver) UpdateBoard(ctx context.Context, input model.UpdateBoardInput) (*model.Board, error) {
	id, err := strconv.ParseUint(input.ID, 10, 64)
	if (err != nil) {
		return nil, err
	}
	err = r.BoardService.UpdateBoard(uint(id), input.Name)
	if (err != nil) {
		return nil, err
	}
	b, err := r.BoardService.GetBoard(uint(id))
	if (err != nil) {
		return nil, err
	}
	return &model.Board{
		ID:        strconv.FormatUint(uint64(b.ID), 10),
		Name:      b.Name,
		CreatedAt: b.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: b.UpdatedAt.Format(utils.TimeFormat),
	}, nil
}

func (r *mutationResolver) DeleteBoard(ctx context.Context, id string) (bool, error) {
	bid, err := strconv.ParseUint(id, 10, 64)
	if (err != nil) {
		return false, err
	}
	err = r.BoardService.DeleteBoard(uint(bid))
	return err == nil, err
}

func (r *queryResolver) Boards(ctx context.Context) ([]*model.Board, error) {
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		return nil, errors.New("未驗證身份")
	}
	boards, err := r.BoardService.GetBoardsByUserID(userID)
	if err != nil {
		return nil, err
	}
	result := make([]*model.Board, 0, len(boards))
	for _, b := range boards {
		result = append(result, &model.Board{
			ID:        strconv.FormatUint(uint64(b.ID), 10),
			Name:      b.Name,
			CreatedAt: b.CreatedAt.Format(utils.TimeFormat),
			UpdatedAt: b.UpdatedAt.Format(utils.TimeFormat),
		})
	}
	return result, nil
}

func (r *queryResolver) Board(ctx context.Context, id string) (*model.Board, error) {
	boardID, err := strconv.ParseUint(id, 10, 64)
	if (err != nil) {
		return nil, err
	}
	b, err := r.BoardService.GetBoard(uint(boardID))
	if (err != nil) {
		return nil, err
	}
	return &model.Board{
		ID:        strconv.FormatUint(uint64(b.ID), 10),
		Name:      b.Name,
		CreatedAt: b.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: b.UpdatedAt.Format(utils.TimeFormat),
	}, nil
}
