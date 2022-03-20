package gomake

import (
	"embed"
	"strings"
	"text/template"

	"git.sr.ht/~ionous/tapestry/dl/spec"
)

//go:embed templates/*
var tempFS embed.FS

func newTemplates(ctx *Context) (*template.Template, error) {
	//
	funcMap := template.FuncMap{
		"Lines": func(s string) []string {
			return strings.Split(s, "\n")
		},
		"Pascal": pascal,
		"Tokenize": func(k string) string {
			return "$" + strings.ToUpper(k)
		},
		//
		// more specific things:
		//
		"ScopeOf": func(typeName string) string {
			return ctx.scopeOf(typeName)
		},
		"IsUnboxed": func(typeName string) (okay bool) {
			_, okay = unbox[typeName]
			return
		},
		"Terms": func(typeSpec *spec.TypeSpec) []Term {
			flow := typeSpec.Spec.Value.(*spec.FlowSpec)
			terms := make([]Term, len(flow.Terms))
			var pubCount int
			for i, t := range flow.Terms {
				pi := -1
				if !t.Private {
					pi = pubCount
					pubCount++
				}
				terms[i] = Term{ctx, flow, t, pi}
			}
			return terms
		},
		"Uses": func(typeSpec *spec.TypeSpec) string {
			return map[string]string{
				spec.UsesSpec_Flow_Opt:  "flow",
				spec.UsesSpec_Slot_Opt:  "slot",
				spec.UsesSpec_Swap_Opt:  "swap",
				spec.UsesSpec_Num_Opt:   "num",
				spec.UsesSpec_Str_Opt:   "str",
				spec.UsesSpec_Group_Opt: "group",
			}[typeSpec.Spec.Choice]
		},

		"RepeatData": func(name string, unboxed bool) any {
			// {{>repeat name=(Pascal name) mod="_Unboxed" el=(Unboxed name)}}
			var mod, el string
			pas := pascal(name)
			if !unboxed {
				el = pas
			} else {
				el = unbox[name]
				mod = "_Unboxed"
			}
			return struct {
				Name, Mod, El string
			}{pas, mod, el}
		},
	}
	return template.New("").Funcs(funcMap).ParseFS(tempFS, "templates/*.tmpl")
}
