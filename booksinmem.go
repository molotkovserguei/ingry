package main

import (
	"context"
	"sync"
)

type booksInMem struct {
	mu    sync.RWMutex
	books map[int]Book
	seq   int
}

var _ BooksDB = (*booksInMem)(nil)

func NewBooksInMem() *booksInMem {
	return &booksInMem{
		books: make(map[int]Book),
	}
}

func (b *booksInMem) List(ctx context.Context) (result []Book, next int, err error) {
	b.mu.RLock()

	response := make([]Book, 0, len(b.books))
	for _, v := range b.books {
		response = append(response, v)
	}

	b.mu.RUnlock()
	return response, 0, nil
}

func (b *booksInMem) Get(ctx context.Context, id int) (Book, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	response, ok := b.books[id]
	if !ok {
		return Book{}, ErrNotFound
	}

	return response, nil
}

func (b *booksInMem) Create(ctx context.Context, preparedBooks ...Book) ([]int, error) {
	ids := make([]int, 0, len(preparedBooks))

	b.mu.Lock()
	for _, newBook := range preparedBooks {
		b.seq++
		newBook.ID = b.seq
		b.books[b.seq] = newBook
		ids = append(ids, b.seq)
	}
	b.mu.Unlock()

	return ids, nil
}

func (b *booksInMem) Update(ctx context.Context, oneBook Book) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.books[oneBook.ID]
	if !ok {
		return ErrNotFound
	}

	b.books[oneBook.ID] = oneBook

	return nil
}

func (b *booksInMem) Delete(ctx context.Context, id int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.books[id]
	if !ok {
		return ErrNotFound
	}

	delete(b.books, id)

	return nil
}
