package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year" gorm:"column:year_issue"`
}

type bookStorage struct {
	booksDB BooksDB
}

func NewBooks(repodb RepositoryDB) *bookStorage {
	return &bookStorage{
		booksDB: repodb.Books(),
	}
}

func (b *bookStorage) List(c echo.Context) error {
	response, _, err := b.booksDB.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("list books: %w", err).Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (b *bookStorage) Get(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.ParseInt(paramID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("id must be int: %s", err.Error()))
	}

	response, err := b.booksDB.Get(c.Request().Context(), int(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.JSON(http.StatusNotFound, fmt.Sprintf("book %s not found", paramID))
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("book %s: %s", paramID, err.Error()))
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

	newIds, err := b.booksDB.Create(c.Request().Context(),
		Book{
			Title:  title,
			Author: author,
			Year:   int(year),
		})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("create book: %w", err).Error())
	}

	return c.String(http.StatusOK, "id:"+strconv.FormatInt(int64(newIds[0]), 10))
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

	oneBook := Book{
		ID:     int(id),
		Title:  c.FormValue("title"),
		Author: c.FormValue("author"),
		Year:   int(year),
	}
	err = b.booksDB.Update(c.Request().Context(), oneBook)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.JSON(http.StatusNotFound, fmt.Sprintf("book %s not found", paramID))
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("book %s: %s", paramID, err.Error()))
	}

	return c.JSON(http.StatusOK, oneBook)
}

func (b *bookStorage) Delete(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.ParseInt(paramID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("id must be int: %s", err.Error()))
	}

	err = b.booksDB.Delete(c.Request().Context(), int(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.JSON(http.StatusNotFound, fmt.Sprintf("book %s not found", paramID))
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("book %s: %s", paramID, err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}
