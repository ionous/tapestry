package doc

import (
	"embed"
	"fmt"
	"html/template"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

//go:embed templates/*
var temFS embed.FS

func docTemplates(g GlobalData) (*template.Template, error) {
	dc := docComments{g}
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
				ret = dc.formatComment(lines)
			}
			return
		},
		"GoComment":  dc.formatComment,
		"LinkByType": g.linkByType,
		"LinkByName": g.linkByName,
		"LinkByIdl":  g.linkByIdl,
		"Testing": func(a, b any) (ret string, err error) {
			err = fmt.Errorf("it worked %v %v", a, b)
			return
		},
	}
	return template.New("").Funcs(funcMap).ParseFS(temFS, "templates/*.tem")
}
