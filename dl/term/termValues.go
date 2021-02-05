package term

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
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
	Init       rt.RecordEval
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

type RecordList struct {
	Name, Kind string
	Init       rt.RecordListEval
}

func (n *Number) String() string {
	return n.Name
}

func (n *Number) Prepare(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetOptionalNumber(run, n.Init, 0); e != nil {
		err = errutil.New("error preparing", n.Name, e)
	} else {
		ret = v
	}
	return
}

func (n *Bool) String() string {
	return n.Name
}

func (n *Bool) Prepare(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetOptionalBool(run, n.Init, false); e != nil {
		err = errutil.New("error preparing", n.Name, e)
	} else {
		ret = v
	}
	return
}

func (n *Text) String() string {
	return n.Name
}

func (n *Text) Prepare(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetOptionalText(run, n.Init, ""); e != nil {
		err = errutil.New("error preparing", n.Name, e)
	} else {
		ret = v
	}
	return
}

func (n *Record) String() string {
	return n.Name
}

func (n *Record) Prepare(run rt.Runtime) (ret g.Value, err error) {
	if v, e := n.prep(run); e != nil {
		err = errutil.New("error preparing", n.Name, e)
	} else {
		ret = v
	}
	return
}

func (n *Record) prep(run rt.Runtime) (ret g.Value, err error) {
	if k, e := run.GetKindByName(n.Kind); e != nil {
		err = e
	} else if init := n.Init; init == nil {
		v := g.RecordOf(k.NewRecord())
		ret = v
	} else if v, e := init.GetRecord(run); e != nil {
		err = e
	} else if want, have := k.Name(), v.Type(); want != have {
		err = errutil.New("expected record of", want, "have", have)
	} else {
		ret = v
	}
	return
}

func (n *Object) String() string {
	return n.Name
}

// Prepare - an object returns a text value ( or error )
// the "type" of the text is set to the kind of the object.
func (n *Object) Prepare(run rt.Runtime) (ret g.Value, err error) {
	if v, e := n.prep(run); e != nil {
		err = errutil.New("error preparing", n.Name, e)
	} else {
		ret = v
	}
	return
}

func (n *Object) prep(run rt.Runtime) (ret g.Value, err error) {
	if k, e := run.GetKindByName(n.Kind); e != nil {
		err = errutil.New("error preparing", n.Name, e)
	} else if v, e := safe.GetOptionalText(run, n.Init, ""); e != nil {
		err = e
	} else if res, e := ConvertName(run, v.String(), k.Name()); e != nil {
		err = e
	} else {
		ret = res
	}
	return
}

func (n *NumList) String() string {
	return n.Name
}

func (n *NumList) Prepare(run rt.Runtime) (ret g.Value, err error) {
	if vs, e := safe.GetOptionalNumbers(run, n.Init, nil); e != nil {
		err = errutil.New("error preparing", n.Name, e)
	} else {
		ret = vs
	}
	return
}

func (n *TextList) String() string {
	return n.Name
}

func (n *TextList) Prepare(run rt.Runtime) (ret g.Value, err error) {
	if vs, e := safe.GetOptionalTexts(run, n.Init, nil); e != nil {
		err = errutil.New("error preparing", n.Name, e)
	} else {
		ret = vs
	}
	return
}

func (n *RecordList) String() string {
	return n.Name
}

func (n *RecordList) Prepare(run rt.Runtime) (ret g.Value, err error) {
	if v, e := n.prep(run); e != nil {
		err = errutil.New("error preparing", n.Name, e)
	} else {
		ret = v
	}
	return
}

func (n *RecordList) prep(run rt.Runtime) (ret g.Value, err error) {
	if k, e := run.GetKindByName(n.Kind); e != nil {
		err = e
	} else if init := n.Init; init == nil {
		v := g.RecordsOf(k.Name(), nil)
		ret = v
	} else if vs, e := init.GetRecordList(run); e != nil {
		err = e
	} else if want, have := k.Name(), vs.Type(); want != have {
		err = errutil.New("expected record of", want, "have", have)
	} else {
		ret = vs
	}
	return
}
