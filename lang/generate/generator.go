package generate

import (
	"errors"
	"io"
	"strings"
	"text/template"

	"git.sr.ht/~ionous/tapestry/lang/compact"
)

// read a "Spec:with group:" message
func ReadGroup(msg compact.Message) (ret Group, err error) {
	var gc groupContent
	if msg.Key != "Spec:with group:" {
		err = errors.New("expected a top level group spec") // ugh.
	} else if n, e := parseString("spec", msg.Args[0], ""); e != nil {
		err = e
	} else if e := readSpec(n, &gc, msg); e != nil {
		err = e
	} else if cmt, e := compact.ExtractComment(msg.Markup); e != nil {
		err = e
	} else {
		ret = Group{n, gc, cmt}
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

// for schemas
type slotList struct {
	slotData
	Types []slotEntry
}

// for non evals it returns the empty string
// ex. address
func (n slotList) ChopEval() (ret string) {
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
	slot := make(map[string]slotList)
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
				slot[op.Name] = slotList{slotData: op}
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
						slot[s] = slotList{Types: []slotEntry{entry}}
					} else {
						a.Types = append(a.Types, entry)
						slot[s] = a
					}
				}
			}
		}
	}
	return q.tmp.ExecuteTemplate(w, "schema.tmpl", struct {
		Name          string
		SchemaComment string
		SchemaId      string
		Flow          []flowData
		Str           []strData
		Num           []numData
		Slot          map[string]slotList
	}{
		"Tell", "A Tapestry story file",
		"https://tapestry.ionous.net/schema/tell/v0",
		flow, str, num, slot,
	})
}
