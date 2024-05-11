// Package doc builds web friendly documentation of the idl.
package doc

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

const ext = ".html"
const SourceUrl = "https://pkg.go.dev/git.sr.ht/~ionous/tapestry@v0.24.4-1"
const typesFolder = "types"
const slotsFolder = "slot"

//go:embed static/*
var staticFS embed.FS

// hack
var allTypes []typeinfo.TypeSet

// outDir can point to a temp directory if need be
func Build(outDir string, idl []typeinfo.TypeSet) (err error) {
	allTypes = idl
	if tem, e := docTemplates(); e != nil {
		err = e
	} else if e := os.MkdirAll(outDir, os.ModePerm); e != nil {
		err = e
	} else if e := CopyStaticFiles(outDir); e != nil {
		err = e
	} else {
		slots := make(SlotMap)
		for _, types := range idl {
			splitBySlot(slots, types)
			os.Mkdir(filepath.Join(outDir, typesFolder), os.ModePerm)
			outFile := filepath.Join(outDir, typesFolder, types.Name+ext)
			if e := Create(outFile, tem, map[string]any{
				"Name":  types.Name,
				"Types": types,
			}); e != nil {
				err = e
				return // early out
			}
		}
		//
		os.Mkdir(filepath.Join(outDir, slotsFolder), os.ModePerm)
		for slot, part := range slots {
			outFile := filepath.Join(outDir, slotsFolder, slot.Name+ext)
			if e := Create(outFile, tem, map[string]any{
				"Name": slot.Name,
				"Slot": slot,
				"Part": part,
			}); e != nil {
				err = e
				return // early out
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

func Create(outFile string, tem *template.Template, data any) (err error) {
	if fp, e := os.Create(outFile); e != nil {
		err = e
	} else if e := tem.ExecuteTemplate(fp, "page.tem", data); e != nil {
		err = fmt.Errorf("%w writing %s", e, outFile)
	}
	return
}

// - split everything into slots
func splitBySlot(out SlotMap, types typeinfo.TypeSet) {
	for _, t := range types.Flow {
		for _, slot := range t.Slots {
			m, ok := out[slot]
			if !ok {
				m = make(PartitionMap)
				out[slot] = m
			}
			// fix:
			pkgName, pkgPath := types.Name, "dl/"+types.Name
			m.addFlow(t, pkgName, pkgPath)
		}
	}
}

type SlotMap map[*typeinfo.Slot]PartitionMap

type PartitionMap map[string]Partition

func (m *PartitionMap) addFlow(t *typeinfo.Flow, pkgName, pkgPath string) {
	if t.Lede != "--" {
		cat := (*m)[t.Lede]
		cat.addFlow(t, pkgName, pkgPath)
		(*m)[t.Lede] = cat
	}
}

type Partition struct {
	Types []Command
	// 1. keep these sorted  by command name  ( the default )
	// alt: sort by label
}

func (cat *Partition) addFlow(t *typeinfo.Flow, pkgName, pkgPath string) {
	cat.Types = append(cat.Types, Command{t, pkgName, pkgPath})
}

type Command struct {
	*typeinfo.Flow
	pkgName, pkgPath string
}

func (c Command) SourceLink() string {
	// fix: need the package name
	name := inflect.Pascal(c.Name)
	return path.Join(SourceUrl, c.pkgPath+"#"+name)
}

func (c Command) Terms() (ret []typeinfo.Term) {
	// filter out private terms;
	for i, t := range c.Flow.Terms {
		if t.Private && ret == nil {
			ret = c.Flow.Terms[:i]
		}
		if !t.Private && ret != nil {
			ret = append(ret, t)
		}
	}
	if ret == nil {
		ret = c.Flow.Terms
	}
	return
}

func TypeLink(t typeinfo.T) (ret string) {
	name := t.TypeName()
	pascal := inflect.Pascal(name)
	switch t.(type) {
	default:
		log.Panicf("unknown type %T", t)
	case *typeinfo.Slot:
		ret = join("../"+slotsFolder+"/", name, ext)

	case *typeinfo.Flow:
		var idl string
	Flow:
		for _, t := range allTypes {
			for _, el := range t.Flow {
				if el.Name == name {
					idl = t.Name
					break Flow
				}
			}
		}
		if len(idl) == 0 {
			log.Panicf("couldnt find flow for %s", name)
		}
		ret = join("../"+typesFolder+"/", idl, ext, "#", pascal)

	case *typeinfo.Str:
		var idl string
	Str:
		for _, t := range allTypes {
			for _, el := range t.Str {
				if el.Name == name {
					idl = t.Name
					break Str
				}
			}
		}
		if len(idl) == 0 {
			log.Panicf("couldnt find str for %s", name)
		}
		ret = join("../"+typesFolder+"/", idl, ext, "#", pascal)

	case *typeinfo.Num:
		var idl string
	Num:
		for _, t := range allTypes {
			for _, el := range t.Num {
				if el.Name == name {
					idl = t.Name
					break Num
				}
			}
		}
		if len(idl) == 0 {
			log.Panicf("couldnt find num for %q", name)
		}
		ret = join("../"+typesFolder+"/", idl, ext, "#", pascal)

	}
	return
}

func join(str ...string) string {
	return strings.Join(str, "")
}
