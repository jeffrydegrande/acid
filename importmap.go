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

var currentCDN CDN
var importMap *ImportMap

type Package struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type ImportMap struct {
	Packages  []Package
	Structure Structure
}

type Structure struct {
	Imports map[string]string `json:"imports,omitempty"`
}

func newImportMap() *ImportMap {
	return &ImportMap{
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
	if importMap == nil {
		importMap = newImportMap()
	}
	p := Package{
		Name: name,
		URL:  url,
	}
	importMap.Packages = append(importMap.Packages, p)
	importMap.Structure.Imports[p.Name] = p.URL
}

func PinAllFrom(fs *embed.FS) {
	CalculateDigests(fs, "static")

	files, err := getFiles(fs, "static")
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

func getFiles(fs *embed.FS, path string) (out []string, err error) {
	if len(path) == 0 {
		path = "."
	}

	entries, err := fs.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fp := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			res, err := getFiles(fs, fp)
			if err != nil {
				return nil, err
			}
			out = append(out, res...)
			continue
		}
		out = append(out, fp)
	}

	return
}

func buildURL(p string, version string) string {
	switch currentCDN {
	case JsDelivr:
		return fmt.Sprintf("https://cdn.jsdelivr.net/npm/%s@%s/+esm", p, version)
	default:
		panic("Unknown CDN")
	}
}

func RenderImportMap() (template.HTML, error) {
	if importMap == nil {
		return "", errors.New("ImportMap hasn't been setup")
	}
	return importMap.Render()
}

func (im *ImportMap) Imports() (template.HTML, error) {
	b, err := json.MarshalIndent(im.Structure, "", "\t")
	if err != nil {
		return "", err
	}
	return template.HTML(b), nil
}

// Render returns a HTML snippet to use in a template
func (im *ImportMap) Render() (template.HTML, error) {
	// TODO: handle modulepreload properly
	t, err := template.New("").Parse(`
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
	err = t.Execute(&buf, im)

	if err != nil {
		return "", err
	}
	return template.HTML(buf.Bytes()), nil
}
