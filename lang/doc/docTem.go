package doc

import (
	"embed"
	"fmt"
	"html/template"
	"strings"

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
		"LinkByType": g.linkByType,
		"LinkByName": g.linkByName,
		"LinkByIdl":  g.linkByIdl,
		"Comment": func(src any) (ret template.HTML, err error) {
			switch src := src.(type) {
			case map[string]any:
				if lines, e := compact.ExtractComment(src); e != nil {
					err = e
				} else if len(lines) > 0 {
					ret = dc.formatComment(lines)
				}
			case []string:
				if len(src) > 0 {
					ret = dc.formatComment(src)
				}
			default:
				err = fmt.Errorf("cant handle comments of %T", src)
			}
			return
		},
		"Snippet": func(src any) (ret string, err error) {
			switch src := src.(type) {
			case map[string]any:
				if lines, e := compact.ExtractComment(src); e != nil {
					err = e
				} else if len(lines) > 0 {
					ret = makeSnippet(lines)
				}
			case []string:
				if len(src) > 0 {
					ret = makeSnippet(src)
				}
			default:
				err = fmt.Errorf("cant handle snippets of %T", src)
			}
			return
		},
	}
	return template.New("").Funcs(funcMap).ParseFS(temFS, "templates/*.tem")
}

func makeSnippet(lines []string) (ret string) {
	str := strings.Join(lines, " ")
	if i := strings.IndexRune(str, '.'); i <= 0 {
		ret = str
	} else {
		ret = str[:i+1]
	}
	return
}
