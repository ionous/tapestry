package list_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/kr/pretty"
)

type visitEach struct {
	visits *[]accum
}

func TestEach(t *testing.T) {
	// primary looping test:
	eachTest(t, []string{
		"Orange", "Lemon", "Mango",
	}, []accum{
		// what we expect to see from the index, first, last, and text values
		// when looping over the list of fruits
		{1, true, false, "Orange"},
		{2, false, false, "Lemon"},
		{3, false, true, "Mango"},
	}, 0 /*... and zero calls to else */)

	// never any middle ground
	eachTest(t, []string{
		"Orange", "Mango",
	}, []accum{
		{1, true, false, "Orange"},
		{2, false, true, "Mango"},
	}, 0 /*... and zero calls to else */)

	// first and last are both true
	eachTest(t, []string{
		"Lime",
	}, []accum{
		{1, true, true, "Lime"},
	}, 0 /*... and zero calls to else */)

	// looping over an empty list
	eachTest(t, nil, nil,
		1 /*... and a single call to else */)
}

func eachTest(t *testing.T, src []string, res []accum, otherwise int) {
	var out []string
	var visits []accum
	each := &list.ListRepeat{
		List: &call.FromTextList{Value: object.Variable("source")},
		As:   ("text"),
		Exe:  []rt.Execute{&visitEach{&visits}},
		Else: &logic.ChooseNothingElse{
			Exe: []rt.Execute{&Write{&out, literal.T("x")}},
		},
	}
	if lt, e := newListTime(src, nil); e != nil {
		t.Fatal(e)
	} else if e := each.Execute(lt); e != nil {
		t.Fatal(src, e)
	} else if d := pretty.Diff(visits, res); len(d) > 0 || len(out) != otherwise {
		t.Fatal(src, out, d)
	}
}

func (v *visitEach) Execute(run rt.Runtime) (err error) {
	if i, e := checkVariable(run, "index", affine.Num); e != nil {
		err = e
	} else if f, e := checkVariable(run, "first", affine.Bool); e != nil {
		err = e
	} else if l, e := checkVariable(run, "last", affine.Bool); e != nil {
		err = e
	} else if t, e := checkVariable(run, "text", affine.Text); e != nil {
		err = e
	} else {
		(*v.visits) = append((*v.visits), accum{
			i.Int(), f.Bool(), l.Bool(), t.String(),
		})
	}
	return
}

type accum struct {
	index int
	first bool
	last  bool
	text  string
}

func checkVariable(run rt.Runtime, name string, aff affine.Affinity) (ret rt.Value, err error) {
	if v, e := run.GetField(meta.Variables, name); e != nil {
		err = e
	} else if e := safe.Check(v, aff); e != nil {
		err = e
	} else {
		ret = v
	}
	return

}
