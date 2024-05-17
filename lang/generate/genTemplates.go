package generate

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"git.sr.ht/~ionous/tapestry/lang/compact"
)

//go:embed templates/*
var tempFS embed.FS

// fix: the template functions for schema probably deserve to be split off
func genTemplates(p *groupSearch) (*template.Template, error) {
	funcMap := template.FuncMap{
		//
		// shared functions:
		//
		"Pascal": Pascal,
		"Title":  Titlecase,
		//
		// schema functions:
		//
		"TrimmedSignature": func(str string) string {
			end := len(str) - 1
			if end > 0 && str[end] == ':' {
				str = str[:end]
			}
			return str
		},
		// a special signature to make refs vs properties in schema clearer.
		// replaces colons and spaces with underscores
		// starts with an underscore to force a separate "namespace"
		// from the plain type names ( which are also lower case )
		"LowerSignature": func(str string) string {
			var out strings.Builder
			prefix := true // start prefixed
			for _, r := range str {
				if prefix {
					out.WriteRune('_')
				}
				if prefix = r == ':'; !prefix {
					if r == ' ' {
						r = '_'
					} else {
						r = unicode.ToLower(r)
					}
					out.WriteRune(r)
				}
			}
			return out.String()
		},
		//
		// go generation functions:
		//
		"CommentLines": compact.ExtractComment,
		"Encode": func(v any) (ret string) {
			return fmt.Sprintf("%#v", v)
		},
		// return the package scope; doesnt care if the tapestry type is exported
		// ( useful for typeinfo references where the go type becomes a primitive value )
		// ( ex. to reference the bool info for a bool type )
		"PackageDot": func(typeName string) (ret string, err error) {
			if t, ok := p.findType(typeName); !ok {
				err = fmt.Errorf("unknown type %q", typeName)
			} else if t.group != p.currentGroup {
				ret = t.group + "."
			}
			return
		},
		// return a scoped go type scoped for the named tapestry type.
		"ScopedType": func(typeName string) (ret string, err error) {
			if t, ok := p.findType(typeName); !ok {
				err = fmt.Errorf("unknown type %q", typeName)
			} else if goName := t.goType(); t.group != p.currentGroup && exported(goName) {
				ret = t.group + "." + goName
			} else {
				ret = goName
			}
			return
		},
		// return a scoped go type for the term, scoped by the package.
		// requires overrides for bool, num, str.
		"TermType": func(term termData) (ret string, err error) {
			if termType := term.Type; term.Private {
				ret = Pascal(termType)
			} else if t, ok := p.findType(termType); !ok {
				err = fmt.Errorf("unknown type %q", termType)
			} else {
				if goName := t.goType(); t.group != p.currentGroup && exported(goName) {
					ret = t.group + "." + goName
				} else {
					ret = goName
				}
				_, isFlow := t.typeData.(flowData)
				if term.Optional && !term.Repeats && isFlow {
					ret = "*" + ret
				}
			}
			if term.Repeats {
				ret = "[]" + ret
			}
			return
		},
	}
	return template.New("").Funcs(funcMap).ParseFS(tempFS, "templates/*.tmpl")
}

// looks at whether the first letter is capitalized
func exported(n string) bool {
	first := n[:1]
	return strings.ToUpper(first) == first
}
