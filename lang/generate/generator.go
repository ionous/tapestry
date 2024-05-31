package generate

import (
	"cmp"
	"errors"
	"io"
	"slices"
	"strings"
	"text/template"

	"git.sr.ht/~ionous/tapestry/lang/compact"
)

// read a "Spec:requires:contains:" message
func ReadSpec(msg compact.Message) (ret Group, err error) {
	if msg.Lede != "spec" {
		err = errors.New("expected a spec") // ugh.
	} else {
		mm := MakeMessageMap(msg)
		if name, e := mm.GetString("", ""); e != nil {
			err = e
		} else if group, e := readSpec(name, mm); e != nil {
			err = e
		} else if comment, e := compact.ExtractComment(msg.Markup); e != nil {
			err = e
		} else {
			ret = Group{name, group, comment}
		}
	}
	return
}

// the generator walks a list of groups,
// producing one group at a time:
// each of which can be written to its own location.
type Generator struct {
	groups *groupSearch
	tmp    *template.Template
	i      int
}

func MakeGenerator(groups []Group) (ret Generator, err error) {
	seeker := newGroupSearch(groups)
	if tmp, e := genTemplates(seeker); e != nil {
		err = e
	} else {
		ret = Generator{groups: seeker, tmp: tmp}
	}
	return
}

// currents the current group
func (q *Generator) group() Group {
	return q.groups.list[q.i-1]
}

// advance before writing
func (q *Generator) Next() (okay bool) {
	if okay = (q.i < len(q.groups.list)); okay {
		q.groups.setCurrent(q.i)
		q.i++
	}
	return
}

// database/sql like interface
type DB interface {
	Write(q string, args ...interface{}) error
}

func (q *Generator) Name() string {
	return q.group().Name
}

// write a go file containing typeinfo for the current group
func (q *Generator) WriteSource(w io.Writer) error {
	return q.writeSource(w, q.group())
}

// record the primary type data for the current group
// using sqlite friendly data to the passed database
func (q *Generator) WriteTable(w DB) error {
	return writeTable(w, q.group())
}

// record the derived type data for the current group
// using sqlite friendly data to the passed database
func (q *Generator) WriteReferences(w DB) error {
	return writeReferences(w, q.group())
}

// for schemas: the slot and the types in that slot.
type slotCommands struct {
	slotData
	cmds []slotEntry
}

// used by the template to walk the slice
func (n slotCommands) Types() []slotEntry {
	slices.SortFunc(n.cmds, func(a, b slotEntry) int {
		return cmp.Compare(a.TypeScope(), b.TypeScope())
	})
	return n.cmds
}

// for non evals it returns the empty string
// ex. address
func (n slotCommands) ChopEval() (ret string) {
	const eval = "_eval"
	if str := n.Name; strings.HasSuffix(str, eval) {
		ret = str[:len(str)-len(eval)]
	}
	return
}

type slotEntry struct {
	Idl, Type string
}

func (n slotEntry) TypeScope() string {
	return n.Idl + "." + n.Type
}

func (q Generator) WriteSchema(w io.Writer) (err error) {
	// tbd, maybe one schema can include others
	var flow []flowData
	var str []strData
	var num []numData
	slot := make(map[string]slotCommands)
	hackForLinks = q.groups

	for q.Next() {
		curr := q.group()
		for _, op := range curr.Str {
			str = append(str, op.(strData))
		}
		for _, op := range curr.Num {
			num = append(num, op.(numData))
		}
		for _, op := range curr.Slot {
			op := op.(slotData)
			if a, ok := slot[op.Name]; !ok {
				slot[op.Name] = slotCommands{slotData: op}
			} else {
				a.slotData = op
				slot[op.Name] = a
			}
		}
		for _, f := range curr.Flow {
			op := f.(flowData)
			if _, private := op.Markup["internal"]; !private {
				flow = append(flow, op)
				entry := slotEntry{Idl: op.Idl, Type: op.Name}
				for _, s := range op.Slots {
					if a, ok := slot[s]; !ok {
						slot[s] = slotCommands{cmds: []slotEntry{entry}}
					} else {
						a.cmds = append(a.cmds, entry)
						slot[s] = a
					}
				}
			}
		}
	}
	return q.tmp.ExecuteTemplate(w, "schema.tem", struct {
		Name          string
		SchemaComment string
		SchemaId      string
		Flow          []flowData
		Str           []strData
		Num           []numData
		Slot          map[string]slotCommands
	}{
		"Tell", "A Tapestry story file",
		"https://tapestry.ionous.net/schema/tell/v0",
		flow, str, num, slot,
	})
}
