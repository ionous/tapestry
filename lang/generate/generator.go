package generate

import (
	"errors"
	"io"
	"text/template"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/support/distill"
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
	} else {
		ret = Group{n, gc}
	}
	return
}

type Generator struct {
	groups groupSearch
	tmp    *template.Template
	err    error
	i      int
}

func MakeGenerator(groups []Group) (ret Generator) {
	seeker := groupSearch{list: groups}
	if tmp, e := genTemplates(&seeker); e != nil {
		ret = Generator{err: e}
	} else {
		ret = Generator{groups: seeker, tmp: tmp}
	}
	return
}

func (q *Generator) Name() string {
	g := q.groups.list[q.i-1]
	return g.Name
}

// advance before writing
func (q *Generator) Next() (okay bool) {
	if okay = q.err != nil || (q.i < len(q.groups.list)); okay {
		q.i++
	}
	return
}

// after collecting groups, write them to targeted location
// ( it tell format, with the tells extension )
func (q *Generator) Write(w io.Writer) (err error) {
	if q.err != nil {
		err = q.err
	} else if at := q.i - 1; at >= len(q.groups.list) {
		err = errors.New("out of range")
	} else {
		g := q.groups.list[at]
		var sigs distill.Registry
		if e := q.write(w, "header", struct {
			Name    string
			Imports []string
		}{g.Name, nil}); e != nil {
			err = e
		} else if e := q.writeContent(w, g.groupContent); e != nil {
			err = e
		} else if e := q.writeTypes(w, "Y_Slot", "*typeinfo.Slot",
			"a list of all slots ( ex. for generating blockly shapes )",
			g.Slot); e != nil {
			err = e
		} else if e := q.writeTypes(w, "Y_Flow", "*typeinfo.Flow",
			"a list of all flows ( ex. for reading blockly blocks )",
			g.Flow); e != nil {
			err = e
		} else if e := q.writeSigs(w, sigs); e != nil {
			err = e
		}
	}
	return
}

// writes a list of typeinfo references
func (q *Generator) writeSigs(w io.Writer, reg distill.Registry) error {
	return q.tmp.ExecuteTemplate(w, "signatures", reg.Sigs)
}

// writes a list of typeinfo references
func (q *Generator) writeTypes(w io.Writer, name, subType, comment string, list []typeData) error {
	data := struct {
		Name, Type, Comment string
		List                []typeData
	}{name, subType, comment, list}
	return q.write(w, "types", data)
}

func (q *Generator) writeContent(w io.Writer, gc groupContent) (err error) {
	for i, cnt := 0, len(gc.Slot); i < cnt && err == nil; i++ {
		n := gc.Slot[i]
		err = q.write(w, "slot", n)
	}
	for i, cnt := 0, len(gc.Flow); i < cnt && err == nil; i++ {
		n := gc.Flow[i]
		err = q.write(w, "flow", n)
	}
	for i, cnt := 0, len(gc.Str); i < cnt && err == nil; i++ {
		n := gc.Str[i].(strData)
		if len(n.Options) > 0 {
			err = q.write(w, "enum", n)
		} else {
			err = q.write(w, "string", n)
		}
	}
	for i, cnt := 0, len(gc.Num); i < cnt && err == nil; i++ {
		n := gc.Num[i]
		err = q.write(w, "num", n)
	}
	return
}

func (q *Generator) write(w io.Writer, name string, data any) error {
	return q.tmp.ExecuteTemplate(w, name+".tmpl", data)
}
