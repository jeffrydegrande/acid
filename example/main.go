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

// TemplRenderer helps setting up an echo renderer. It's not perfect, but it
// does the traick and lets you call c.Render to return from your echo handlers
type TemplRenderer struct{}

// Render implements echo.Renderer. It's not perfect because the 2nd string
// argument is entirely unnecessary, but it's required by echo.Renderer. We can
// just ignore it though.
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

	// we can tie a FS route to acid. The interface for this isn't perfect but works.
	// check assets/embed.go to check how things tie together.
	e.GET("/assets/*", echo.WrapHandler(
		http.StripPrefix("/assets/", http.FileServer(acid.FS()))))

	e.Logger.Fatal(e.Start(":6969"))
}
