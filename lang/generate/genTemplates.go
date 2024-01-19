package generate

import (
	"embed"
	"text/template"

	"git.sr.ht/~ionous/tapestry/support/distill"
)

//go:embed templates/*
var tempFS embed.FS

var Pascal = distill.Pascal

func genTemplates(p TypeFinder) (*template.Template, error) {
	funcMap := template.FuncMap{
		"Pascal": distill.Pascal,
		"Scoped": func(typeName string) (ret string, err error) {
			if n, ok := p.findScope(typeName); !ok {
				// err = fmt.Errorf("unknown type %q", termType)
				ret = typeName // temp"
			} else if len(n) > 0 {
				ret = n
			}
			return
		},
		// return a scoped go type for the term
		// requires overrides for bool, num, str.
		"TermType": func(t termData) (ret string, err error) {
			if termType := t.Type; t.Private {
				ret = Pascal(termType)
			} else if typeName := p.findType(termType); len(typeName) > 0 {
				ret = typeName
			} else {
				// err = fmt.Errorf("unknown type %q", termType)
				ret = "bool // temp"
			}
			return
		},
	}
	return template.New("").Funcs(funcMap).ParseFS(tempFS, "templates/*.tmpl")
}

type TypeFinder interface {
	// given a lowercase type name, find the go type
	findType(n string) string
	// only for flow and slot
	// given a lowercase type name, find the package
	findScope(n string) (string, bool)
}
