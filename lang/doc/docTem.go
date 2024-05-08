package doc

import (
	"embed"
	"html/template"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

//go:embed templates/*
var tempFS embed.FS

type PageInfo struct {
	typeinfo.TypeSet
}

func docTemplates() (*template.Template, error) {
	funcMap := template.FuncMap{
		"Pascal":     inflect.Pascal,
		"Capitalize": inflect.Capitalize,
		// "Encode": func(v any) (ret string) {
		// 	return fmt.Sprintf("%#v", v)
		// },
		// // return the package scope; doesnt care if the tapestry type is exported
		// // ( useful for typeinfo references where the go type becomes a primitive value )
		// // ( ex. to reference the bool info for a bool type )
		// "PackageDot": func(typeName string) (ret string, err error) {
		// 	if t, ok := p.findType(typeName); !ok {
		// 		err = fmt.Errorf("unknown type %q", typeName)
		// 	} else if t.group != p.currentGroup {
		// 		ret = t.group + "."
		// 	}
		// 	return
		// },
		// // return a scoped go type scoped for the named tapestry type.
		// "ScopedType": func(typeName string) (ret string, err error) {
		// 	if t, ok := p.findType(typeName); !ok {
		// 		err = fmt.Errorf("unknown type %q", typeName)
		// 	} else if goName := t.goType(); t.group != p.currentGroup && exported(goName) {
		// 		ret = t.group + "." + goName
		// 	} else {
		// 		ret = goName
		// 	}
		// 	return
		// },
		// // return a scoped go type for the term, scoped by the package.
		// // requires overrides for bool, num, str.
		// "TermType": func(term termData) (ret string, err error) {
		// 	if termType := term.Type; term.Private {
		// 		ret = Pascal(termType)
		// 	} else if t, ok := p.findType(termType); !ok {
		// 		err = fmt.Errorf("unknown type %q", termType)
		// 	} else {
		// 		if goName := t.goType(); t.group != p.currentGroup && exported(goName) {
		// 			ret = t.group + "." + goName
		// 		} else {
		// 			ret = goName
		// 		}
		// 		_, isFlow := t.typeData.(flowData)
		// 		if term.Optional && !term.Repeats && isFlow {
		// 			ret = "*" + ret
		// 		}
		// 	}
		// 	if term.Repeats {
		// 		ret = "[]" + ret
		// 	}
		// 	return
		// },
	}
	return template.New("").Funcs(funcMap).ParseFS(tempFS, "templates/*.tem")
}
