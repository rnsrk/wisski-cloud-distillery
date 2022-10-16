package static

import (
	"encoding/json"
	"html/template"
)

// Assets represents a group of assets to be included inside a template.
//
// Assets are generated using the 'build.mjs' script.
// The script is called using 'go:generate', which stores variables in the form of 'Assets{{Name}}' inside this package.
//
// The build script roughly works as follows:
// - Delete any previously generated distribution directory.
// - Bundle the entrypoint sources under 'src/entry/{{Name}}/index.{ts,css}' together with the base './src/base/index.{ts,css}'
// - Store the output inside the 'dist' directory
// - Generate new constants of the form {{Name}}
//
// Each asset group should be registered as a parameter to the 'go:generate' line.
type Assets struct {
	Scripts string // <script> tags inserted by the asset
	Styles  string // <link> tags inserted by the asset
}

//go:generate node build.mjs HomeHome ControlIndex ControlInstance

// MustParse parses a new template from the given source
// and registers the Asset functions to it.
// See [Assets.RegisterFuncs].
func (assets *Assets) MustParse(value string) *template.Template {
	return template.Must(assets.RegisterFuncs(template.New("")).Parse(value))
}

// RegisterFuncs registers three new template functions called "JS", "CSS" and "json".
//
// "JS" and "CSS" take no arguments, and return appropriate tags to be inserted into html.
// json takes a single argument of any type, and returns it's encoding as a string to be inserted into the page.
func (assets *Assets) RegisterFuncs(t *template.Template) *template.Template {
	return t.Funcs(template.FuncMap{
		"JS":  func() template.HTML { return template.HTML(assets.Scripts) },
		"CSS": func() template.HTML { return template.HTML(assets.Styles) },
		"json": func(data any) (string, error) {
			bytes, err := json.Marshal(data)
			return string(bytes), err
		},
	})
}
