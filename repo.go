package main

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type RepositoryDB interface {
	Books() BooksDB
}

type BooksDB interface {
	Get(ctx context.Context, id int) (Book, error)
	List(ctx context.Context) (result []Book, next int, err error)
	Create(ctx context.Context, preparedBooks ...Book) ([]int, error)
	Update(ctx context.Context, oneBook Book) error
	Delete(ctx context.Context, id int) error
}
