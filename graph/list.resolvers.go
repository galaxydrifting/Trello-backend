package graph

import (
	"context"
	"errors"
	"strconv"
	"trello-backend/graph/model"
	"trello-backend/pkg/utils"

	"github.com/graph-gophers/dataloader"
)

// List 相關 resolver function

func (r *mutationResolver) CreateList(ctx context.Context, input model.CreateListInput) (*model.List, error) {
	bid, err := strconv.ParseUint(input.BoardID, 10, 64)
	if err != nil {
		return nil, err
	}
	l, err := r.ListService.CreateList(uint(bid), input.Name)
	if err != nil {
		return nil, err
	}
	return &model.List{
		ID:        strconv.FormatUint(uint64(l.ID), 10),
		Name:      l.Name,
		BoardID:   strconv.FormatUint(uint64(l.BoardID), 10),
		CreatedAt: l.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: l.UpdatedAt.Format(utils.TimeFormat),
		Position:  int32(l.Position),
	}, nil
}

func (r *mutationResolver) UpdateList(ctx context.Context, input model.UpdateListInput) (*model.List, error) {
	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	err = r.ListService.UpdateList(uint(id), input.Name)
	if err != nil {
		return nil, err
	}
	l, err := r.ListService.GetListByID(uint(id))
	if err != nil {
		return nil, err
	}
	return &model.List{
		ID:        strconv.FormatUint(uint64(l.ID), 10),
		Name:      l.Name,
		BoardID:   strconv.FormatUint(uint64(l.BoardID), 10),
		CreatedAt: l.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: l.UpdatedAt.Format(utils.TimeFormat),
		Position:  int32(l.Position),
	}, nil
}

func (r *mutationResolver) DeleteList(ctx context.Context, id string) (bool, error) {
	lid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return false, err
	}
	err = r.ListService.DeleteList(uint(lid))
	return err == nil, err
}

func (r *mutationResolver) MoveList(ctx context.Context, input model.MoveListInput) (*model.List, error) {
	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	err = r.ListService.MoveList(uint(id), int(input.NewPosition))
	if err != nil {
		return nil, err
	}
	l, err := r.ListService.GetListByID(uint(id))
	if err != nil {
		return nil, err
	}
	return &model.List{
		ID:        strconv.FormatUint(uint64(l.ID), 10),
		Name:      l.Name,
		BoardID:   strconv.FormatUint(uint64(l.BoardID), 10),
		CreatedAt: l.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: l.UpdatedAt.Format(utils.TimeFormat),
		Position:  int32(l.Position),
	}, nil
}

func (r *queryResolver) Lists(ctx context.Context, boardID string) ([]*model.List, error) {
	bid, err := strconv.ParseUint(boardID, 10, 64)
	if err != nil {
		return nil, err
	}
	lists, err := r.ListService.GetLists(uint(bid))
	if err != nil {
		return nil, err
	}
	result := make([]*model.List, 0, len(lists))
	for _, l := range lists {
		result = append(result, &model.List{
			ID:        strconv.FormatUint(uint64(l.ID), 10),
			Name:      l.Name,
			BoardID:   strconv.FormatUint(uint64(l.BoardID), 10),
			CreatedAt: l.CreatedAt.Format(utils.TimeFormat),
			UpdatedAt: l.UpdatedAt.Format(utils.TimeFormat),
			Position:  int32(l.Position),
		})
	}
	return result, nil
}

func (r *queryResolver) List(ctx context.Context, id string) (*model.List, error) {
	listID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	l, err := r.ListService.GetListByID(uint(listID))
	if err != nil {
		return nil, err
	}
	return &model.List{
		ID:        strconv.FormatUint(uint64(l.ID), 10),
		Name:      l.Name,
		BoardID:   strconv.FormatUint(uint64(l.BoardID), 10),
		CreatedAt: l.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: l.UpdatedAt.Format(utils.TimeFormat),
		Position:  int32(l.Position),
	}, nil
}

// Cards is the resolver for the cards field.
func (r *listResolver) Cards(ctx context.Context, obj *model.List) ([]*model.Card, error) {
	loaders := For(ctx)
	if loaders == nil {
		return nil, errors.New("dataloader not found in context")
	}
	thunk := loaders.CardsByListID.Load(ctx, dataloader.StringKey(obj.ID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	cards, ok := result.([]*model.Card)
	if !ok {
		return nil, errors.New("unexpected dataloader result type")
	}
	return cards, nil
}
