package main

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type booksPSQL struct {
	db *gorm.DB
}

var _ BooksDB = (*booksPSQL)(nil)

func NewBooksPSQL(db *gorm.DB) *booksPSQL {
	return &booksPSQL{
		db: db,
	}
}

func (b *booksPSQL) Get(ctx context.Context, id int) (Book, error) {
	book := Book{}
	r := b.db.First(&book, id)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return book, ErrNotFound
		}
		return book, r.Error
	}

	return book, nil
}

func (b *booksPSQL) List(ctx context.Context) (result []Book, next int, err error) {
	var books []Book
	r := b.db.Find(&books)
	if r.Error != nil {
		return nil, 0, r.Error
	}

	return books, 0, nil
}

func (b *booksPSQL) Create(ctx context.Context, preparedBooks ...Book) ([]int, error) {
	ids := make([]int, 0, len(preparedBooks))

	result := b.db.Create(&preparedBooks)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < int64(len(preparedBooks)) {
		return nil, fmt.Errorf("%d of %d was created", result.RowsAffected, len(preparedBooks))
	}

	for _, v := range preparedBooks {
		ids = append(ids, v.ID)
	}

	return ids, nil
}

func (b *booksPSQL) Update(ctx context.Context, oneBook Book) error {
	result := b.db.Save(&oneBook)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records was updated")
	}

	return nil
}

func (b *booksPSQL) Delete(ctx context.Context, id int) error {
	result := b.db.Delete(&Book{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records deleted")
	}

	return nil
}
