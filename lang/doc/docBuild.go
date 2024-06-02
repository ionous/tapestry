// Package doc builds web friendly documentation of the idl.
package doc

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

const SourceUrl = "https://pkg.go.dev/git.sr.ht/~ionous/tapestry"
const typesFolder = "idl"
const slotFolder = "slot"
const baseUrl = "/api"

//go:embed static/*
var staticFS embed.FS

// outDir can point to a temp directory if need be
func Build(outDir string, idl []typeinfo.TypeSet) (err error) {
	g := makeGlobalData(idl)
	if tem, e := docTemplates(g); e != nil {
		err = e
	} else if e := os.MkdirAll(outDir, os.ModePerm); e != nil {
		err = e
	} else if e := CopyStaticFiles(outDir); e != nil {
		err = e
	} else {
		apiHome := filepath.Join(outDir)
		if e := Create(apiHome, tem, map[string]any{
			"Name":      "Reference",
			"SourceUrl": SourceUrl,
			"Slots":     g.slots,
			"AllTypes":  g.types,
		}); e != nil {
			err = e
			return // early out
		}

		// generate an alphabetical index of all commands
		alphaList := filepath.Join(outDir, "alpha")
		if e := Create(alphaList, tem, map[string]any{
			"Name":     "Index",
			"Commands": g.allCommands,
			"Slots":    g.slots,
			"Num":      g.num,
			"Str":      g.str,
		}); e != nil {
			err = e
		}

		//
		for _, types := range idl {
			var cmds []FlowInfo
			for _, t := range g.allCommands {
				if t.Idl == types.Name {
					cmds = append(cmds, t)
				}
			}
			// generate idl files:
			outFile := filepath.Join(outDir, typesFolder, types.Name)
			if e := Create(outFile, tem, map[string]any{
				"Name":     types.Name,
				"Types":    types,
				"Commands": cmds,
			}); e != nil {
				err = e
				return // early out
			}
		}
		// generate slot docs:
		for _, slot := range g.slots {
			var cmds []FlowInfo
			for _, t := range g.allCommands {
				if slices.Contains(t.Slots, slot.Slot) {
					cmds = append(cmds, t)
				}
			}
			outFile := filepath.Join(outDir, slotFolder, slot.Name)
			if e := Create(outFile, tem, map[string]any{
				"Name":     slot.Name,
				"Slot":     slot,
				"Commands": cmds,
			}); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func CopyStaticFiles(outDir string) (err error) {
	for _, fileName := range []string{"style.css", "custom.css"} {
		if b, e := fs.ReadFile(staticFS, "static/"+fileName); e != nil {
			err = e
		} else if fp, e := os.Create(filepath.Join(outDir, fileName)); e != nil {
			err = e
		} else {
			_, err = fp.Write(b)
		}
	}
	return
}

func Create(outPath string, tem *template.Template, data any) (err error) {
	os.MkdirAll(outPath, os.ModePerm)
	outFile := filepath.Join(outPath, "index.html")
	if fp, e := os.Create(outFile); e != nil {
		err = e
	} else if e := tem.ExecuteTemplate(fp, "page.tem", data); e != nil {
		err = fmt.Errorf("%w writing %s", e, outFile)
	}
	return
}
