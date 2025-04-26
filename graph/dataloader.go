package graph

import (
	"context"
	"strconv"
	"trello-backend/graph/model"
	"trello-backend/internal/services"

	"trello-backend/pkg/utils"

	"github.com/graph-gophers/dataloader"
)

type Loaders struct {
	CardsByListID *dataloader.Loader
}

// CardsBatchFn 批次查詢多個 List 的 Cards
func CardsBatchFn(cardService services.CardService) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		results := make([]*dataloader.Result, len(keys))
		listIDs := make([]uint, len(keys))
		for i, k := range keys {
			id, _ := strconv.ParseUint(k.String(), 10, 64)
			listIDs[i] = uint(id)
		}
		// 一次查詢所有 listIDs 的 cards
		cardsMap, err := cardService.GetCardsByListIDs(listIDs)
		for i, k := range keys {
			id, _ := strconv.ParseUint(k.String(), 10, 64)
			cards := cardsMap[uint(id)]
			modelCards := make([]*model.Card, 0, len(cards))
			for _, c := range cards {
				modelCards = append(modelCards, &model.Card{
					ID:        strconv.FormatUint(uint64(c.ID), 10),
					Title:     c.Title,
					Content:   strToPtr(c.Content),
					ListID:    strconv.FormatUint(uint64(c.ListID), 10),
					CreatedAt: c.CreatedAt.Format(utils.TimeFormat),
					UpdatedAt: c.UpdatedAt.Format(utils.TimeFormat),
					Position:  int32(c.Position),
				})
			}
			results[i] = &dataloader.Result{Data: modelCards, Error: err}
		}
		return results
	}
}

// context key
var loadersKey = &struct{}{}

// Middleware: 將 dataloader 注入 context
func DataloaderMiddleware(cardService services.CardService) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		loaders := &Loaders{
			CardsByListID: dataloader.NewBatchedLoader(CardsBatchFn(cardService)),
		}
		return context.WithValue(ctx, loadersKey, loaders)
	}
}

// 從 context 取得 loaders
func For(ctx context.Context) *Loaders {
	val := ctx.Value(loadersKey)
	if val == nil {
		return nil
	}
	return val.(*Loaders)
}
