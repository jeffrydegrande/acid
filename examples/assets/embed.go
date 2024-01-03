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

// We embed our entire asset directory. acid will work with the embed.FS so you
// don't really need to handle any asset stuff on your build server. I just
// commit anything generated to git and have my CI / build server pick up
// nothing but source coe.

//go:embed static/*
var Assets embed.FS

func init() {
	// I'm only supporting jsdelivr for now. I'll add more CDNs as I need them.
	// This statement configures how to rewrite urls for the pins we'll define below.
	acid.UseCDN(acid.JsDelivr)

	// This is how you add a dependency. These will be added to the importmap
	acid.Pin("htmx.org", "1.9.10")

	// this will grab all assets in the Assets defined above and wire up digesting
	// for these files.
	acid.PinAllFrom(&Assets)
}
