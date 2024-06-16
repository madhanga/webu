package webu

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed web
var web embed.FS

func createTmplCache() (map[string]*template.Template, error) {
	var cache = map[string]*template.Template{}

	pages, err := fs.Glob(web, "web/*.page.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		var tmplName = filepath.Base(page)
		var tmplNameStart = strings.Split(tmplName, ".")[0]

		tmpl, err := template.New(tmplName).ParseFS(web, page)
		if err != nil {
			return nil, err
		}

		for _, t := range tmpl.Templates() {
			if t.Name() != tmplName && t.Tree.ParseName == tmplName {
				cache[fmt.Sprintf("%s.%s", tmplNameStart, t.Name())] = t
			}
		}

		matches, err := fs.Glob(web, "web/*.layout.html")
		if len(matches) > 0 {
			tmpl, err = tmpl.ParseFS(web, "web/*.layout.html")
			if err != nil {
				return nil, err
			}
		}

		cache[tmplNameStart] = tmpl
	}

	return cache, nil
}
