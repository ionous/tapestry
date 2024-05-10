// Package doc builds web friendly documentation of the idl.
package doc

import (
	"embed"
	"html/template"
	"io"
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

//go:embed static/*
var staticFS embed.FS

// outDir can point to a temp directory if need be
func Build(outDir string, idl []typeinfo.TypeSet) (err error) {
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
		}
		//
		for slot, part := range slots {
			// fix: maybe some directories instead leading underscore?
			// ( something is needed to avoid collision with str/num types )
			outFile := filepath.Join(outDir, "slot_"+slot.Name+ext)
			if fp, e := os.Create(outFile); e != nil {
				err = e
				break
			} else if e := SlotPage(fp, tem, slot, part); e != nil {
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

func SlotPage(w io.Writer, tem *template.Template, slot *typeinfo.Slot, part PartitionMap) error {
	return tem.ExecuteTemplate(w, "page.tem", struct {
		Slot *typeinfo.Slot
		Part PartitionMap
	}{
		slot, part,
	})
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
	case *typeinfo.Flow:
		// FIX!
		// and are there types that dont implement a slot?
		// probably

	case *typeinfo.Slot:
		ret = join("slot_", name, ext)

	case *typeinfo.Str:
		ret = join("type_str", ext, "#", pascal)
	case *typeinfo.Num:
		// hrm: extension links
		ret = join("type_num", ext, "#", pascal)
	default:
		log.Panicf("unknown type %T", t)
	}
	return
}

func join(str ...string) string {
	return strings.Join(str, "")
}
