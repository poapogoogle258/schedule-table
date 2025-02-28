package repository

import "gorm.io/gorm"

type Transaction interface {
	Begin() *gorm.DB
}

type transaction struct {
	db *gorm.DB
}

func (t *transaction) Begin() *gorm.DB {
	return t.db.Begin()
}

func NewTransaction(db *gorm.DB) Transaction {
	return &transaction{db: db}
}
