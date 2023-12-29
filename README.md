# acid

Acid is a library for doing #NoBuild frontends in Go projects.

[![Go Report Card](https://goreportcard.com/badge/github.com/jeffrydegrande/acid)](https://goreportcard.com/report/github.com/jeffrydegrande/acid)

## Features

- #NoBuild strategy for frontend assets
- Leans heavily into embed.FS: all your asset files are compiled into your binary so you only have to deploy a single file.
- Generates importmaps for your JS dependencies. currently only supporting JsDelivr as a CDN
- Automatically generates and maps digests for busting caches.

In my setup I use the tailwindcss binary to prerender some CSS because I'm rather bad at CSS. You don't have to do that.
You can add your vanilla css and it will be picked up.

## Usage

To install `acid`, run:

```sh
go get github.com/jeffrydegrande/acid
```

Include `acid` in your project. Here's an example of how I'm using it right now:

```go
package assets

import (
	"embed"
	"github.com/jeffrydegrande/acid"
)

// Include tailwindcss, this will end up in the Assets FS
//go:generate tailwindcss build -o ./static/css/tailwind.css --minify

//go:embed static/*
var Assets embed.FS

func init() {
	acid.UseCDN(acid.JsDelivr)
	acid.Pin("htmx.org", "1.9.10")
	acid.Pin("flatpickr", "4.6.9")
	acid.Pin("@milkdown/core", "7.3.2")
	acid.Pin("@milkdown/ctx", "7.3.2")
	acid.Pin("@milkdown/plugin-listener", "7.3.2")
	acid.Pin("@milkdown/preset-commonmark", "7.3.2")
	acid.PinAllFrom(&Assets)
}
```

Three interesting things are going on here.

1. I do use `tailwindcss` and I use go:generate to produce my output before
   go:embed pix up the result
2. embed.FS pulls _everything_ in the static/ directory into something I can
   later read from memory. This means everything, all my css, all my images,
   all my js will be compiled into my binary.
3. I'm leaning on the init function to "load everything". This means I do need
   to import this module at an approriate time. I do it in the module that sets
   up my Echo server with "import \_ "github.com/jeffrydegrande/<project>/assets"
   in this case.

## Configuration

TODO:

## License

`acid` is made available under the MIT License. See the LICENSE file for more details.
