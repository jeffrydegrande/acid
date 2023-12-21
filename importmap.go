package acid

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
)

type CDN int

const (
	JsDelivr CDN = iota
)

var currentCDN CDN
var importMap *ImportMap

func UseCDN(cdn CDN) {
	currentCDN = cdn
}

func Pin(p string, version string) {
	if importMap == nil {
		importMap = NewImportMap([]Package{})
	}

	importMap.Packages = append(importMap.Packages, Package{
		Name: p,
		URL:  buildURL(p, version),
	})
}

func buildURL(p string, version string) string {
	switch currentCDN {
	case JsDelivr:
		return fmt.Sprintf("https://cdn.jsdelivr.net/npm/%s@%s/+esm", p, version)
	default:
		panic("Unknown CDN")
	}
}

type Package struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type ImportMap struct {
	Packages  []Package
	Structure structure
}

type structure struct {
	Imports map[string]string `json:"imports,omitempty"`
}

func RenderImportMap() (template.HTML, error) {
	if importMap == nil {
		return "", errors.New("ImportMap hasn't been setup")
	}
	return importMap.Render()
}

func NewImportMap(packages []Package) *ImportMap {
	im := &ImportMap{
		Packages: packages,
		Structure: structure{
			Imports: make(map[string]string),
		},
	}

	// prepare structure for rendering
	for _, entry := range im.Packages {
		im.Structure.Imports[entry.Name] = entry.URL
	}

	return im
}

// NOTE:  we're depending on behavior in the golang specification here by assuming that
// init functions in the same package are run in alphabetical order of the module name.
/*
func init() {
	files, err := GetAssets()
	if err != nil {
		// panicing here isn't great because we're not in a situation
		// where we can pass the error up the stack
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file) != ".js" {
			continue
		}

		modulePath := AssetsWithDigests.ReverseMap[file]
		moduleName := strings.TrimSuffix(strings.TrimPrefix(file, "static/javascript/"), ".js")

		importMap.Structure.Imports[moduleName] = fmt.Sprintf("/assets/%s", modulePath)
	}
}
*/

func (im *ImportMap) Imports() (template.HTML, error) {
	b, err := json.MarshalIndent(im.Structure, "", "\t")
	if err != nil {
		return "", err
	}
	return template.HTML(b), nil
}

// Render returns a HTML snippet to use in a template
func (im *ImportMap) Render() (template.HTML, error) {

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

	return template.HTML(buf.String()), nil
}
