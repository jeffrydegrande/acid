package main

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/jeffrydegrande/acid"
	_ "github.com/jeffrydegrande/acid/examples/assets"
	"github.com/jeffrydegrande/acid/examples/views"
	"github.com/labstack/echo/v4"
)

var (
	errTemplateNotFound = errors.New("Template not found")
)

type TemplRenderer struct{}

func (t TemplRenderer) Render(w io.Writer, _ string, data interface{}, c echo.Context) error {

	if templData, ok := data.(templ.Component); ok {
		return templData.Render(context.Background(), w)

	}
	return errTemplateNotFound
}

func main() {
	e := echo.New()
	e.Renderer = TemplRenderer{}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "", views.Demo("Hello, World!"))
	})
	e.GET("/assets/*", echo.WrapHandler(
		http.StripPrefix("/assets/", http.FileServer(acid.FS()))))

	e.Logger.Fatal(e.Start(":6969"))
}
