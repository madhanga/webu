package webu

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

type Template struct {
	cache map[string]*template.Template
}

func LoadTemplates(ui embed.FS) (*Template, error) {
	var cache = map[string]*template.Template{}

	pages, err := fs.Glob(ui, "ui/**/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		var tmplName = filepath.Base(page)
		var tmplNameStart = strings.Split(tmplName, ".")[0]

		tmpl, err := template.New(tmplName).ParseFS(ui, page)
		if err != nil {
			return nil, err
		}

		for _, t := range tmpl.Templates() {
			if t.Name() == "css" || t.Name() == "js" || t.Name() == "main" {
				continue
			}
			if t.Name() != tmplName && t.Tree.ParseName == tmplName {
				cache[fmt.Sprintf("%s.%s", tmplNameStart, t.Name())] = t
			}
		}

		tmpl, err = tmpl.ParseFS(ui, "ui/index.html")
		if err != nil {
			return nil, err
		}

		cache[tmplNameStart] = tmpl
	}

	return &Template{cache: cache}, nil
}

func (t *Template) Render(w http.ResponseWriter, tmplName string, data any) error {
	var tmpl *template.Template
	var found bool
	if tmpl, found = t.cache[tmplName]; !found {
		return errors.New("template " + tmplName + " found in cache")
	}
	return tmpl.Execute(w, data)
}

func Err(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf(err.Error())
}
