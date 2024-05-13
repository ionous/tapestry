package doc

import (
	"embed"
	"html/template"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

//go:embed templates/*
var temFS embed.FS

func docTemplates(d *docComments) (*template.Template, error) {
	funcMap := template.FuncMap{
		"Pascal":     inflect.Pascal,
		"Capitalize": inflect.Capitalize,
		"Title":      inflect.Titlecase,
		"Camel":      inflect.Camelize,
		// "Lines":      extractLines,
		"MarkupComment": func(m map[string]any) (ret template.HTML, err error) {
			if lines, e := compact.ExtractComment(m); e != nil {
				err = e
			} else if len(lines) > 0 {
				ret = d.formatComment(lines)
			}
			return
		},
		"GoComment": d.formatComment,
		"TypeLink":  TypeLink,
	}
	return template.New("").Funcs(funcMap).ParseFS(temFS, "templates/*.tem")
}
