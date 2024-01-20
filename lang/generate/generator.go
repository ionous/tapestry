package generate

import (
	"bytes"
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
	groups *groupSearch
	tmp    *template.Template
	i      int
}

func MakeGenerator(groups []Group) (ret Generator, err error) {
	seeker := &groupSearch{list: groups}
	if tmp, e := genTemplates(seeker); e != nil {
		err = e
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
	if okay = (q.i < len(q.groups.list)); okay {
		q.i++
	}
	return
}

// after collecting groups, write them to targeted location
// ( it tell format, with the tells extension )
func (q *Generator) Write(w io.Writer) (err error) {
	if at := q.i - 1; at >= len(q.groups.list) {
		err = errors.New("out of range")
	} else {
		g := q.groups.setCurrent(at)
		// this uses the generation to determine inter package references
		// so we have to build the body first, and later
		// write the import header above it.
		var body bytes.Buffer
		var sigs distill.Registry
		if e := q.writeContent(&body, g.groupContent); e != nil {
			err = e
		} else if e := q.writeTypes(&body, "slot", "( ex. for generating blockly shapes )",
			g.Slot); e != nil {
			err = e
		} else if e := q.writeTypes(&body, "flow", "( ex. for reading blockly blocks )",
			g.Flow); e != nil {
			err = e
		} else if e := q.writeSigs(&body, sigs); e != nil {
			err = e
		}
		dl := "git.sr.ht/~ionous/tapestry/dl/" // fix: customize?
		ti := "git.sr.ht/~ionous/tapestry/lang/typeinfo"
		imports := q.groups.getRefs(dl, ti)
		if e := q.write(w, "header", struct {
			Name    string
			Imports []string
		}{g.Name, imports}); e != nil {
			err = e
		} else {
			w.Write(body.Bytes())
		}
	}
	return
}

// writes a list of typeinfo references
func (q *Generator) writeSigs(w io.Writer, reg distill.Registry) error {
	return q.write(w, "signatures", reg.Sigs)
}

// writes a list of typeinfo references
func (q *Generator) writeTypes(w io.Writer, name, comment string, list []typeData) (err error) {
	if len(list) > 0 {
		data := struct {
			Name, Comment string
			List          []typeData
		}{name, comment, list}
		err = q.write(w, "types", data)
	}
	return
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
