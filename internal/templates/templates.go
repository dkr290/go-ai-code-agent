package templates

import (
	"html/template"
	"path/filepath"
)

// LoadTemplates parses all HTML files in the specified pattern
func LoadTemplates(pattern string) (*template.Template, error) {
	tmpl := template.New("") // Create an empty template set
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		_, err = tmpl.ParseFiles(f)
		if err != nil {
			return nil, err
		}
	}
	return tmpl, nil
}
