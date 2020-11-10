package term

import (
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type Number struct {
	Name string        // parameter name
	Init rt.NumberEval // default value
}
type Bool struct {
	Name string
	Init rt.BoolEval
}
type Text struct {
	Name string
	Init rt.TextEval
}
type Record struct {
	Name, Kind string
}
type Object struct {
	Name, Kind string
	Init       rt.TextEval
}
type NumList struct {
	Name string
	Init rt.NumListEval
}
type TextList struct {
	Name string
	Init rt.TextListEval
}

// type RecordList struct {
// 	Name, Kind string
// 	// possibly with an initial size generating a zero list.
// }

func (n *Number) String() string {
	return n.Name
}

func (n *Number) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalNumber(run, n.Init, 0); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, g.FloatOf(v))
	}
	return
}

func (n *Bool) String() string {
	return n.Name
}

func (n *Bool) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalBool(run, n.Init, false); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, g.BoolOf(v))
	}
	return
}

func (n *Text) String() string {
	return n.Name
}

func (n *Text) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalText(run, n.Init, ""); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, g.StringOf(v))
	}
	return
}

func (n *Record) String() string {
	return n.Name
}

func (n *Record) Prepare(run rt.Runtime, p *Terms) (err error) {
	if k, e := run.GetKindByName(n.Kind); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, g.RecordOf(k.NewRecord()))
	}
	return
}

func (n *Object) String() string {
	return n.Name
}

func (n *Object) Prepare(run rt.Runtime, p *Terms) (err error) {
	if v, e := rt.GetOptionalText(run, n.Init, ""); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, g.StringOf(v))
	}
	return
}

func (n *NumList) String() string {
	return n.Name
}

func (n *NumList) Prepare(run rt.Runtime, p *Terms) (err error) {
	if vs, e := rt.GetOptionalNumbers(run, n.Init, nil); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, g.FloatsOf(vs))
	}
	return
}

func (n *TextList) String() string {
	return n.Name
}

func (n *TextList) Prepare(run rt.Runtime, p *Terms) (err error) {
	if vs, e := rt.GetOptionalTexts(run, n.Init, nil); e != nil {
		err = e
	} else {
		p.AddTerm(n.Name, g.StringsOf(vs))
	}
	return
}

// func (n *RecordList) String() string {
// 	return n.Name
// }

// func (n *RecordList) Prepare(run rt.Runtime, p *Terms) (err error) {
// 	p.AddTerm(n.Name, g.NewNewRecordSlice(n.Kind))
// 	return
// }
