package main

import (
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	// mydb, _ := NewRepoInMem()
	mydb, err := NewRepoDB()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	books := NewBooks(mydb)
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			return nil
		}
	})

	e.GET("/books", books.List)
	e.GET("/books/", books.List)
	e.GET("/books/:id", books.Get)
	e.POST("/books", books.Create)
	e.POST("/books/", books.Create)
	e.PUT("/books/:id", books.Update)
	e.DELETE("/books/:id", books.Delete)

	e.Logger.Fatal(e.Start(":1323"))
}
