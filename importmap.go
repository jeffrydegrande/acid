package acid

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

type CDN int

const (
	JsDelivr CDN = iota
)

var (
	currentCDN CDN
	im         *importMap

	errImportMapNotSetup = errors.New("importmap hasn't been setup")
)

type Package struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type importMap struct {
	Packages  []Package
	Structure Structure
}

type Structure struct {
	Imports map[string]string `json:"imports,omitempty"`
}

func newImportMap() *importMap {
	return &importMap{
		Packages: []Package{},
		Structure: Structure{
			Imports: make(map[string]string),
		},
	}
}

func UseCDN(cdn CDN) {
	currentCDN = cdn
}

func Pin(name string, version string) {
	pin(name, buildURL(name, version))
}

func pin(name, url string) {
	if im == nil {
		im = newImportMap()
	}

	pkg := Package{
		Name: name,
		URL:  url,
	}

	im.Packages = append(im.Packages, pkg)
	im.Structure.Imports[pkg.Name] = pkg.URL
}

func PinAllFrom(fs *embed.FS) {
	err := CalculateDigests(fs, "static")
	if err != nil {
		panic(err)
	}

	files, err := listFiles(fs, "static")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file) != ".js" {
			continue
		}

		url := assetsWithDigests.ReverseMap[file]
		name := strings.TrimSuffix(strings.TrimPrefix(file, "static/javascript/"), ".js")
		pin(name, fmt.Sprintf("/assets/%s", url))
	}
}

func listFiles(fileSystem *embed.FS, path string) ([]string, error) {
	if len(path) == 0 {
		path = "."
	}

	entries, err := fileSystem.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// var out []string
	out := make([]string, 0, len(entries))

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			res, err := listFiles(fileSystem, fullPath)
			if err != nil {
				return nil, err
			}

			out = append(out, res...)

			continue
		}

		out = append(out, fullPath)
	}

	return out, err
}

func buildURL(p string, version string) string {
	switch currentCDN {
	case JsDelivr:
		return fmt.Sprintf("https://cdn.jsdelivr.net/npm/%s@%s/+esm", p, version)
	default:
		panic("Unknown CDN")
	}
}

func Packages() []Package {
	return im.Packages
}

func renderImportMap() (template.HTML, error) {
	if im == nil {
		return "", errImportMapNotSetup
	}

	return im.Render()
}

func (m *importMap) Imports() (template.HTML, error) {
	b, err := json.MarshalIndent(m.Structure, "", "\t")
	if err != nil {
		return "", err
	}

	return template.HTML(b), nil // #nosec G203 -- Don't want this to be escaped
}

// Render returns a HTML snippet to use in a template.
func (m *importMap) Render() (template.HTML, error) {
	tmpl, err := template.New("").Parse(`
<script type="importmap">
	{{ .Imports }}
</script>

{{ range .Packages }}
		<link rel="modulepreload" href="{{ .URL }}">
{{ end }}
`)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, m)

	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil // #nosec G203 -- Don't want this to be escaped
}
