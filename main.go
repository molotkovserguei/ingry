package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	books := NewBooks()
	e := echo.New()

	e.GET("/books", books.List)
	e.GET("/books/", books.List)
	e.GET("/books/:id", books.Get)
	e.POST("/books", books.Create)
	e.POST("/books/", books.Create)
	e.PUT("/books/:id", books.Update)
	e.DELETE("/books/:id", books.Delete)

	e.Logger.Fatal(e.Start(":1323"))
}
