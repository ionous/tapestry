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
			lines := strings.Split(s, "\n")
			// fix... backwards compat:
			// if last := len(lines) - 1; last >= 0 && strings.HasSuffix(lines[last], ".") {
			// 	lines[last] = strings.TrimSuffix(lines[last], ".")
			// }
			return lines
		},
		"Pascal": pascal,
		//
		// more specific things:
		//
		"ScopeOf": func(typeName string) string {
			return ctx.scopeOf(typeName)
		},
		"IsUnboxed": func(typeName string) (okay bool) {
			_, okay = ctx.unbox[typeName]
			return
		},
		"Terms": func(block *spec.TypeSpec) []Term {
			return ctx.TermsOf(block)
		},
		"Uses": func(block *spec.TypeSpec) string {
			return specShortName(block)
		},
		"RepeatData": func(name string, unboxed bool) any {
			// {{>repeat name=(Pascal name) mod="_Unboxed" el=(Unboxed name)}}
			var mod, el string
			pas := pascal(name)
			if !unboxed {
				el = pas
			} else {
				el = ctx.unbox[name]
				mod = "_Unboxed"
			}
			return struct {
				Name, Mod, El string
			}{pas, mod, el}
		},
	}
	return template.New("").Funcs(funcMap).ParseFS(tempFS, "templates/*.tmpl")
}

func specShortName(block *spec.TypeSpec) string {
	return map[string]string{
		spec.UsesSpec_Flow_Opt:  "flow",
		spec.UsesSpec_Slot_Opt:  "slot",
		spec.UsesSpec_Swap_Opt:  "swap",
		spec.UsesSpec_Num_Opt:   "num",
		spec.UsesSpec_Str_Opt:   "str",
		spec.UsesSpec_Group_Opt: "group",
	}[block.Spec.Choice]
}
