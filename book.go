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
		return &echo.HTTPError{
			Internal: err,
			Message:  "list books",
			Code:     http.StatusInternalServerError,
		}
	}

	return c.JSON(http.StatusOK, response)
}

func (b *bookStorage) Get(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.ParseInt(paramID, 10, 32)
	if err != nil {
		return idMustBeInt(err)
	}

	response, err := b.booksDB.Get(c.Request().Context(), int(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return bookNotFound(paramID)
		}
		return bookDBError(err, paramID)
	}

	return c.JSON(http.StatusOK, response)
}

func (b *bookStorage) Create(c echo.Context) error {
	title := c.FormValue("title")
	author := c.FormValue("author")
	sYear := c.FormValue("year")

	year, err := strconv.ParseInt(sYear, 10, 32)
	if err != nil {
		return yearMustBeInt(err)
	}

	newIds, err := b.booksDB.Create(c.Request().Context(),
		Book{
			Title:  title,
			Author: author,
			Year:   int(year),
		})
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "create book",
			Code:     http.StatusInternalServerError,
		}
	}

	return c.String(http.StatusOK, "id:"+strconv.FormatInt(int64(newIds[0]), 10))
}

func (b *bookStorage) Update(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.ParseInt(paramID, 10, 32)
	if err != nil {
		return idMustBeInt(err)
	}

	sYear := c.FormValue("year")
	year, err := strconv.ParseInt(sYear, 10, 32)
	if err != nil {
		return yearMustBeInt(err)
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
			return bookNotFound(paramID)
		}
		return bookDBError(err, paramID)
	}

	return c.JSON(http.StatusOK, oneBook)
}

func (b *bookStorage) Delete(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.ParseInt(paramID, 10, 32)
	if err != nil {
		return idMustBeInt(err)
	}

	err = b.booksDB.Delete(c.Request().Context(), int(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return bookNotFound(paramID)
		}
		return bookDBError(err, paramID)
	}

	return c.NoContent(http.StatusNoContent)
}

func idMustBeInt(err error) *echo.HTTPError {
	return &echo.HTTPError{
		Internal: err,
		Message:  "id must be int",
		Code:     http.StatusBadRequest,
	}
}

func yearMustBeInt(err error) *echo.HTTPError {
	return &echo.HTTPError{
		Internal: err,
		Message:  "year must be int",
		Code:     http.StatusBadRequest,
	}
}

func bookNotFound(id string) *echo.HTTPError {
	return &echo.HTTPError{
		Internal: ErrNotFound,
		Message:  fmt.Errorf("book %s not found", id),
		Code:     http.StatusNotFound,
	}
}

func bookDBError(err error, id string) *echo.HTTPError {
	return &echo.HTTPError{
		Internal: err,
		Message:  fmt.Sprintf("book %s", id),
		Code:     http.StatusInternalServerError,
	}
}
