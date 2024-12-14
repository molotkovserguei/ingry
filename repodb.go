package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type repoEntities struct {
	books *booksPSQL
}

func (r *repoEntities) Books() BooksDB {
	return r.books
}

func NewRepoDB() (*repoEntities, error) {
	dsn := "host=localhost user=psq_user password=psq_password dbname=ebs port=5433 sslmode=disable"

	ns := schema.NamingStrategy{
		TablePrefix:   "ing.", // schema name
		SingularTable: false,
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: ns,
	})
	if err != nil {
		return nil, err
	}

	return &repoEntities{
		books: NewBooksPSQL(db),
	}, nil
}
