package templates

import (
	"html/template"
	"os"
	"path/filepath"
)

func Parse(root string) (*template.Template, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"truncate": Truncate,
	})

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err := tmpl.ParseFiles(path)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
