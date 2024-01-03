package assets

import (
	"embed"

	"github.com/jeffrydegrande/acid"
)

// Including tailwindcss in your project isn't exactly #NoBuild,
// but you can totally do it. I use this in my projects because I
// am not got at all at CSS. Running `go generate` will generate
// what you need but you do need to have a tailwind.config.js file
// pointing at your templ files. See assets/tailwind.config.js for
// an example
//go:generate tailwindcss build -o ./static/css/tailwind.css --minify

//go:embed static/*
var Assets embed.FS

func init() {
	acid.UseCDN(acid.JsDelivr)
	acid.Pin("htmx.org", "1.9.10")
	acid.PinAllFrom(&Assets)
}
