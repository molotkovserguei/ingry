package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

type book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type bookStorage struct {
	mu    sync.RWMutex
	books map[int]book
	seq   int
}

func NewBooks() *bookStorage {
	return &bookStorage{
		books: make(map[int]book),
	}
}

func (b *bookStorage) List(c echo.Context) error {
	b.mu.RLock()

	response := make([]book, 0, len(b.books))
	for _, v := range b.books {
		response = append(response, v)
	}

	b.mu.RUnlock()
	return c.JSON(http.StatusOK, response)
}

func (b *bookStorage) Get(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.ParseInt(paramID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("id must be int: %s", err.Error()))
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	response, ok := b.books[int(id)]
	if !ok {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("book %s not found", paramID))
	}

	return c.JSON(http.StatusOK, response)
}

func (b *bookStorage) Create(c echo.Context) error {
	title := c.FormValue("title")
	author := c.FormValue("author")
	sYear := c.FormValue("year")

	year, err := strconv.ParseInt(sYear, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("year must be int: %s", err.Error()))
	}

	b.mu.Lock()
	b.seq++
	newBook := book{
		ID:     b.seq,
		Title:  title,
		Author: author,
		Year:   int(year),
	}
	b.books[b.seq] = newBook
	b.mu.Unlock()

	return c.String(http.StatusOK, "id:"+strconv.FormatInt(int64(newBook.ID), 10))
}

func (b *bookStorage) Update(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.ParseInt(paramID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("id must be int: %s", err.Error()))
	}

	sYear := c.FormValue("year")
	year, err := strconv.ParseInt(sYear, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("year must be int: %s", err.Error()))
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.books[int(id)]
	if !ok {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("book %s not found", paramID))
	}

	oneBook := book{
		ID:     int(id),
		Title:  c.FormValue("title"),
		Author: c.FormValue("author"),
		Year:   int(year),
	}
	b.books[oneBook.ID] = oneBook

	return c.JSON(http.StatusOK, oneBook)
}

func (b *bookStorage) Delete(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.ParseInt(paramID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("id must be int: %s", err.Error()))
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.books[int(id)]
	if !ok {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("book %s not found", paramID))
	}

	delete(b.books, int(id))

	return c.NoContent(http.StatusNoContent)
}
