// Package doc builds web friendly documentation of the idl.
package doc

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"slices"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

const SourceUrl = "https://pkg.go.dev/git.sr.ht/~ionous/tapestry"
const typesFolder = "idl"
const slotFolder = "slot"
const baseUrl = "/api"

// outDir can point to a temp directory if need be
func Build(outDir string, idl []typeinfo.TypeSet) (err error) {
	g := makeGlobalData(idl)
	if tem, e := docTemplates(g); e != nil {
		err = e
	} else if e := os.MkdirAll(outDir, os.ModePerm); e != nil {
		err = e
	} else if e := BuildApiIndex(outDir, tem, g); e != nil {
		err = e
	} else if e := BuildIdl(outDir, tem, g); e != nil {
		err = e
	} else if e := BuildSlots(outDir, tem, g); e != nil {
		err = e
	}
	return
}

func BuildApiIndex(outDir string, tem *template.Template, g GlobalData) (err error) {
	// api root file
	apiIndex := filepath.Join(outDir, "_index")
	return WriteDocFile(apiIndex, tem, map[string]any{
		"Name":      "Reference",
		"SourceUrl": SourceUrl,
		"Slots":     g.slots,
		"AllTypes":  g.types,
	})
}

// func BuildAlpha((tem*template.Template, g GlobalData, outDir string)) ( err error) {
// generate an alphabetical index of all commands
// alphaList := filepath.Join(outDir, "alpha")
// if e := WriteDocFile(alphaList, tem, map[string]any{
// 	"Name":     "Index",
// 	"Commands": g.allCommands,
// 	"Slots":    g.slots,
// 	"Num":      g.num,
// 	"Str":      g.str,
// }); e != nil {
// 	err = e
// }

// }

func BuildIdl(outDir string, tem *template.Template, g GlobalData) (err error) {
	subDir := filepath.Join(outDir, typesFolder)
	if e := WriteIndex(subDir, tem, "idl"); e != nil {
		err = e
	} else {
		for _, types := range g.types {
			var cmds []FlowInfo
			for _, t := range g.allCommands {
				if t.Idl == types.Name {
					cmds = append(cmds, t)
				}
			}
			// generate idl files:
			outFile := filepath.Join(subDir, types.Name)
			if e := WriteDocFile(outFile, tem, map[string]any{
				"Name":           types.Name,
				"Types":          types,
				"Commands":       cmds,
				"HasPublicSlots": hasPublicSlots(types),
			}); e != nil {
				err = e
				break
			}
		}
	}
	return
}
func BuildSlots(outDir string, tem *template.Template, g GlobalData) (err error) {
	subDir := filepath.Join(outDir, slotFolder)
	if e := WriteIndex(subDir, tem, "slot"); e != nil {
		err = e
	} else { // generate slot docs:
		for _, slot := range g.slots {
			if internal, _ := slot.Markup["internal"].(bool); !internal {
				var cmds []FlowInfo
				for _, t := range g.allCommands {
					if slices.Contains(t.Slots, slot.Slot) {
						cmds = append(cmds, t)
					}
				}
				outFile := filepath.Join(subDir, slot.Name)
				if e := WriteDocFile(outFile, tem, map[string]any{
					"Name":     slot.Name,
					"Slot":     slot,
					"Commands": cmds,
				}); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func hasPublicSlots(ts typeinfo.TypeSet) (okay bool) {
	for _, slot := range ts.Slot {
		if ok, _ := slot.Markup["internal"].(bool); !ok {
			okay = true
			break
		}
	}
	return
}

func WriteDocFile(outPath string, tem *template.Template, data any) (err error) {
	outFile := filepath.Join(outPath + ".html")
	if fp, e := os.Create(outFile); e != nil {
		err = e
	} else if e := tem.ExecuteTemplate(fp, "page.tem", data); e != nil {
		err = fmt.Errorf("%w writing %s", e, outFile)
	}
	return
}

func WriteIndex(outPath string, tem *template.Template, title string) (err error) {
	os.MkdirAll(outPath, os.ModePerm)
	outFile := filepath.Join(outPath, "_index.html")
	if fp, e := os.Create(outFile); e != nil {
		err = e
	} else if e := tem.ExecuteTemplate(fp, "index.tem", map[string]any{
		"Title": title,
	}); e != nil {
		err = fmt.Errorf("%w writing %s", e, outFile)
	}
	return
}
