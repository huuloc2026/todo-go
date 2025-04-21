package model

import "github.com/huuloc2026/go-to-do.git/common"

type TodoItem struct {
	common.SQLModel
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status"`
}

type ToDoItemCreation struct {
	Id          int    `json:"id" gorm:"column:id;"`
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status"`
}

func (ToDoItemCreation) TableName() string {
	return "todoTables"
}

type ToDoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title;"`
	Description *string `json:"description" gorm:"column:description;"`
	Status      *string `json:"status" gorm:"column:status"`
}

func (TodoItem) TableName() string {
	return "todoTables"
}
func (ToDoItemUpdate) TableName() string {
	return "todoTables"
}
