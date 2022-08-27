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
		//
		// more specific things:
		//
		"ScopeOf": func(typeName string) string {
			return ctx.scopeOf(typeName)
		},
		"Unbox": func(typeName string) string {
			return ctx.unbox[typeName]
		},
		"Terms": func(block *spec.TypeSpec) []Term {
			return ctx.TermsOf(block)
		},
		"UserComment": func(block *spec.TypeSpec) (ret []string) {
			switch cmt := block.Markup["comment"].(type) {
			case string:
				ret = []string{cmt}
			case []string:
				ret = cmt
			}
			return
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
