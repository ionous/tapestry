package gomake

import (
	"embed"
	"strings"
	"text/template"
	"unicode"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
)

//go:embed templates/*
var tempFS embed.FS

type Context struct {
	currentGroup string
	types        rs.TypeSpecs
}

func (ctx *Context) GroupOf(typeName string) (ret string, okay bool) {
	if types, ok := ctx.types.Types[typeName]; ok {
		if len(types.Groups) > 0 {
			ret = types.Groups[0]
			okay = true
		}
	}
	return
}

// private types dont have a typespec
// ( ex. those declared manually by an implementation )
func (ctx *Context) GetTypeSpec(typeName string) (ret *spec.TypeSpec, okay bool) {
	typeSpec, ok := ctx.types.Types[typeName]
	return typeSpec, ok
}

func newTemplates(ctx *Context) (*template.Template, error) {
	// when generating some kinds of simple types...
	// replace the specified typename with specified primitive.
	// types that map to numbers, etc. are added as unbox automatically.
	unbox := map[string]string{"text": "string", "bool": "bool"}
	// underscore_name to PascalCase
	pascal := func(s string) string {
		var out strings.Builder
		for _, p := range strings.Split(strings.ToLower(s), "_") {
			for i, q := range p {
				out.WriteRune(unicode.ToUpper(q))
				out.WriteString(p[i+1:])
				break
			}
		}
		return out.String()
	}
	// return the package qualifier for the passed typename
	// ( doesnt add the typename to the return value )
	scopeOf := func(typeName string) (ret string) {
		if g, ok := ctx.GroupOf(typeName); ok && g != ctx.currentGroup {
			ret = g + "."
		}
		return
	}
	// return the fully qualified typename
	scopedName := func(typeName string, ignoreUnboxing bool) (ret string) {
		if unboxType, ok := unbox[typeName]; ok && !ignoreUnboxing {
			ret = unboxType
		} else {
			ret = scopeOf(typeName) + pascal(typeName)
		}
		return
	}
	//
	funcMap := template.FuncMap{
		"Lines": func(s string) []string {
			return strings.Split(s, "\n")
		},
		"Pascal": pascal,
		// more specific things:
		// while having to produce the type could be done in template, it appears in several places.
		"TermType": func(term *spec.TermSpec) (ret string) {
			var typeName, qualifier string
			if t := term.Type; len(t) > 0 {
				typeName = t
			} else if n := term.Name; len(n) > 0 {
				typeName = n
			} else {
				typeName = term.Key
			}
			termType, ok := ctx.GetTypeSpec(typeName) // the referenced type.
			if term.Repeats {
				qualifier = "[]"
			} else if term.Optional && ok && termType.Spec.Choice == spec.UsesSpec_Flow_Opt {
				qualifier = "*" // pointer to a flow
			}
			return qualifier + scopedName(typeName, false)
		},
		"IsUnboxed": func(typeSpec *spec.TypeSpec) (okay bool) {
			_, okay = unbox[typeSpec.Name]
			return
		},
		"TermIsAnon": func(typeSpec *spec.TypeSpec, term spec.TermSpec) (okay bool) {
			flow := typeSpec.Spec.Value.(*spec.FlowSpec)
			if flow.Trim {
				for _, t := range flow.Terms {
					if t.Private {
						continue // only analyze public terms
					} else if t == term {
						okay = true
					}
					break
				}
			}
			return
		},
	}
	return template.New("").Funcs(funcMap).ParseFS(tempFS, "templates/*.tmpl")
}

func includes(strs []string, str string) (ret bool) {
	for _, el := range strs {
		if el == str {
			ret = true
			break
		}
	}
	return
}
