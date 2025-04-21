package storage

import (
	"context"

	"github.com/huuloc2026/go-to-do.git/modules/items/model"
)

func (s *sqlType) CreateItem(context context.Context, data *model.ToDoItemCreation) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}
