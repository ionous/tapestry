package generate

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed templates/*
var tempFS embed.FS

func genTemplates(p TypeFinder) (*template.Template, error) {
	funcMap := template.FuncMap{
		"Pascal": Pascal,
		"Encode": func(v any) (ret string) {
			return fmt.Sprintf("%#v", v)
		},
		// return the package scope; doesnt care if the tapestry type is exported
		// ( useful for typeinfo references where the go type becomes a primitive value )
		// ( ex. to reference the bool info for a bool type )
		"PackageDot": func(typeName string) (ret string, err error) {
			if pack, goName := p.findType(typeName); len(goName) == 0 {
				err = fmt.Errorf("unknown type %q", typeName)
			} else if len(pack) > 0 {
				ret = pack + "."
			}
			return
		},
		// return a scoped go type scoped for the named tapestry type.
		"ScopedType": func(typeName string) (ret string, err error) {
			if pack, goName := p.findType(typeName); len(goName) == 0 {
				err = fmt.Errorf("unknown type %q", typeName)
			} else if len(pack) > 0 && exported(goName) {
				ret = pack + "." + goName
			} else {
				ret = goName
			}
			return
		},
		// return a scoped go type for the term, scoped by the package.
		// requires overrides for bool, num, str.
		"TermType": func(t termData) (ret string, err error) {
			if termType := t.Type; t.Private {
				ret = Pascal(termType)
			} else if pack, goName := p.findType(termType); len(goName) == 0 {
				err = fmt.Errorf("unknown type %q", termType)
			} else if len(pack) > 0 && exported(goName) {
				ret = pack + "." + goName
			} else {
				ret = goName
			}
			if t.Repeats {
				ret = "[]" + ret
			}
			return
		},
	}
	return template.New("").Funcs(funcMap).ParseFS(tempFS, "templates/*.tmpl")
}

// looks at whether the the
func exported(n string) bool {
	first := n[:1]
	return strings.ToUpper(first) == first
}

type TypeFinder interface {
	// given a lowercase type name, find the go package and type
	findType(n string) (string, string)
}
