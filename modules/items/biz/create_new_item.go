package biz

import (
	"context"

	"github.com/huuloc2026/go-to-do.git/modules/items/model"
)

type CreateItemStorage interface {
	CreateItem(ctx context.Context, data *model.ToDoItemCreation) error
}
type createNewItemBiz struct {
	store CreateItemStorage
}

func NewCreateItemBiz(store CreateItemStorage) *createNewItemBiz {
	return &createNewItemBiz{store}
}
