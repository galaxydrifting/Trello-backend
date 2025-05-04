package graph

import (
	"context"
	"strconv"
	"trello-backend/graph/model"
	"trello-backend/pkg/utils"
)

// Card 相關 resolver function

func (r *mutationResolver) CreateCard(ctx context.Context, input model.CreateCardInput) (*model.Card, error) {
	lid, err := strconv.ParseUint(input.ListID, 10, 64)
	if err != nil {
		return nil, err
	}
	boardID := uint(0)
	if input.BoardID != "" {
		bid, err := strconv.ParseUint(input.BoardID, 10, 64)
		if err != nil {
			return nil, err
		}
		boardID = uint(bid)
	}
	c, err := r.CardService.CreateCard(uint(lid), boardID, input.Title, ptrToStr(input.Content))
	if err != nil {
		return nil, err
	}
	return &model.Card{
		ID:        strconv.FormatUint(uint64(c.ID), 10),
		Title:     c.Title,
		Content:   strToPtr(c.Content),
		ListID:    strconv.FormatUint(uint64(c.ListID), 10),
		CreatedAt: c.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: c.UpdatedAt.Format(utils.TimeFormat),
		Position:  int32(c.Position),
	}, nil
}

func (r *mutationResolver) UpdateCard(ctx context.Context, input model.UpdateCardInput) (*model.Card, error) {
	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	err = r.CardService.UpdateCard(uint(id), input.Title, ptrToStr(input.Content))
	if err != nil {
		return nil, err
	}
	c, err := r.CardService.GetCardByID(uint(id))
	if err != nil {
		return nil, err
	}
	return &model.Card{
		ID:        strconv.FormatUint(uint64(c.ID), 10),
		Title:     c.Title,
		Content:   strToPtr(c.Content),
		ListID:    strconv.FormatUint(uint64(c.ListID), 10),
		CreatedAt: c.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: c.UpdatedAt.Format(utils.TimeFormat),
		Position:  int32(c.Position),
	}, nil
}

func (r *mutationResolver) DeleteCard(ctx context.Context, id string) (bool, error) {
	cid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return false, err
	}
	err = r.CardService.DeleteCard(uint(cid))
	return err == nil, err
}

func (r *mutationResolver) MoveCard(ctx context.Context, input model.MoveCardInput) (*model.Card, error) {
	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	targetListID, err := strconv.ParseUint(input.TargetListID, 10, 64)
	if err != nil {
		return nil, err
	}
	err = r.CardService.MoveCard(uint(id), uint(targetListID), int(input.NewPosition))
	if err != nil {
		return nil, err
	}
	c, err := r.CardService.GetCardByID(uint(id))
	if err != nil {
		return nil, err
	}
	return &model.Card{
		ID:        strconv.FormatUint(uint64(c.ID), 10),
		Title:     c.Title,
		Content:   strToPtr(c.Content),
		ListID:    strconv.FormatUint(uint64(c.ListID), 10),
		CreatedAt: c.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: c.UpdatedAt.Format(utils.TimeFormat),
		Position:  int32(c.Position),
	}, nil
}

func (r *queryResolver) Cards(ctx context.Context, listID string) ([]*model.Card, error) {
	lid, err := strconv.ParseUint(listID, 10, 64)
	if err != nil {
		return nil, err
	}
	cards, err := r.CardService.GetCards(uint(lid))
	if err != nil {
		return nil, err
	}
	result := make([]*model.Card, 0, len(cards))
	for _, c := range cards {
		result = append(result, &model.Card{
			ID:        strconv.FormatUint(uint64(c.ID), 10),
			Title:     c.Title,
			Content:   strToPtr(c.Content),
			ListID:    strconv.FormatUint(uint64(c.ListID), 10),
			CreatedAt: c.CreatedAt.Format(utils.TimeFormat),
			UpdatedAt: c.UpdatedAt.Format(utils.TimeFormat),
			Position:  int32(c.Position),
		})
	}
	return result, nil
}

func (r *queryResolver) Card(ctx context.Context, id string) (*model.Card, error) {
	cid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	c, err := r.CardService.GetCardByID(uint(cid))
	if err != nil {
		return nil, err
	}
	return &model.Card{
		ID:        strconv.FormatUint(uint64(c.ID), 10),
		Title:     c.Title,
		Content:   strToPtr(c.Content),
		ListID:    strconv.FormatUint(uint64(c.ListID), 10),
		CreatedAt: c.CreatedAt.Format(utils.TimeFormat),
		UpdatedAt: c.UpdatedAt.Format(utils.TimeFormat),
		Position:  int32(c.Position),
	}, nil
}
