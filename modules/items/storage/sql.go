package storage

import "gorm.io/gorm"

type sqlType struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *sqlType {
	return &sqlType{db: db}
}
