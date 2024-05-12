package generate

import (
	"errors"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"io"
	"text/template"
)

// read a "Spec:with group:" message
func ReadSpec(msg compact.Message) (ret Group, err error) {
	var gc groupContent
	if msg.Key != "Spec:with group:" {
		err = errors.New("expected a top level group spec") // ugh.
	} else if n, e := parseString("spec", msg.Args[0], ""); e != nil {
		err = e
	} else if e := readSpec(&gc, msg); e != nil {
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
