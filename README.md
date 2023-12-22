# acid

Acid is a library for doing #NoBuild frontends in Go projects.

[![Go Report Card](https://goreportcard.com/badge/github.com/jeffrydegrande/acid)](https://goreportcard.com/report/github.com/jeffrydegrande/acid)

## Features

- No-build strategy for frontend assets
- Easy management of dependencies and static files
- Works great with e.g. tailwindcss

## Installation

To install `acid`, run:

```sh
go get github.com/jeffrydegrande/acid
```

## Usage

Include `acid` in your Go project to manage frontend assets and dependencies. Here's an example of how to use it:

```go
package assets

import (
  "embed"
  "github.com/jeffrydegrande/acid"
)

// Embedding static assets and views
//go:embed static/\*
var Assets embed.FS

//go:embed views/layouts/_ views/templates/_
var Views embed.FS

func init() {
  acid.UseCDN(acid.JsDelivr)
  acid.Pin("@hotwired/stimulus", "3.2.1")
  // ... additional dependencies ...
  acid.CalculateDigests(Assets, "static")
}
```

## Configuration

TODO:

## License

`acid` is made available under the MIT License. See the LICENSE file for more details.
