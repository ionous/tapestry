package generate

import (
	"bytes"
	"io"
)

// after collecting groups, write them to targeted location
// ( it tell format, with the tells extension )
func (q *Generator) writeSource(w io.Writer, g Group) (err error) {
	// this uses the generation to determine inter package references
	// so we have to build the body first, and later
	// write the import header above it.
	var body bytes.Buffer
	if e := q.writeContent(&body, g.groupContent); e != nil {
		err = e
	} else {
		dl := "git.sr.ht/~ionous/tapestry/dl/" // fix: customize?
		var extras []string
		for _, str := range g.Str {
			if str := str.(strData); len(str.Options) > 0 {
				extras = append(extras, "strconv")
				break
			}
		}
		extras = append(extras, "git.sr.ht/~ionous/tapestry/lang/typeinfo")
		imports := q.groups.getImports(dl, extras...)
		if e := q.write(w, "header", struct {
			Name    string
			Imports []string
		}{g.Name, imports}); e != nil {
			err = e
		} else {
			w.Write(body.Bytes())
			err = q.writeFooter(w,
				g.Name, g.Reg,
				typeList{
					Type:    "slot",
					Comment: "( ex. for generating blockly shapes )",
					List:    g.Slot,
				},
				typeList{
					Type:    "flow",
					Comment: "( ex. for reading blockly blocks )",
					List:    g.Flow,
				},
				typeList{
					Type: "str",
					List: g.Str,
				},
				typeList{
					Type: "num",
					List: g.Num,
				},
			)
		}
	}
	return
}

// writes a list of typeinfo references
func (q *Generator) writeFooter(w io.Writer, name string, reg Registry, types ...typeList) error {
	reg.Sort()
	return q.write(w, "footer", struct {
		Name       string
		Types      []typeList
		Signatures []Signature
	}{name, types, reg})
}

type typeList struct {
	Type    string
	Comment string
	List    []typeData
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
