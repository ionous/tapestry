// Package doc builds web friendly documentation of the idl.
package doc

import (
	"html/template"
	"io"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

const ext = ".html"

// outDir can point to a temp directory if need be
func Build(outDir string, idl []typeinfo.TypeSet) (err error) {
	if e := os.MkdirAll(outDir, os.ModePerm); e != nil {
		err = e
	} else if tem, e := docTemplates(); e != nil {
		err = e
	} else {
		slots := make(SlotMap)
		for _, types := range idl {
			splitBySlot(slots, types)
		}
		//
		for slot, part := range slots {
			outFile := filepath.Join(outDir, slot.Name+ext)
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

func SlotPage(w io.Writer, tem *template.Template, slot *typeinfo.Slot, part PartitionMap) error {
	return tem.ExecuteTemplate(w, "page.tem", struct {
		Slot *typeinfo.Slot
		Part PartitionMap
	}{
		slot, part,
	})
}

func groupByLede(types typeinfo.TypeSet) PartitionMap {
	out := make(PartitionMap)
	for _, t := range types.Flow {
		out.addFlow(t)
	}
	return out
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
			m.addFlow(t)
		}
	}
}

type SlotMap map[*typeinfo.Slot]PartitionMap

type PartitionMap map[string]Partition

func (m *PartitionMap) addFlow(t *typeinfo.Flow) {
	if t.Lede != "--" {
		cat := (*m)[t.Lede]
		cat.addFlow(t)
		(*m)[t.Lede] = cat
	}
}

type Partition struct {
	Types []Command
	// 1. keep these sorted  by command name  ( the default )
	// alt: sort by label
}

func (cat *Partition) addFlow(t *typeinfo.Flow) {
	cat.Types = append(cat.Types, Command{t})
}

type Command struct {
	*typeinfo.Flow
}
