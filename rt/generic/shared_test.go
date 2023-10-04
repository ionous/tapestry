package generic_test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
)

func (n *recordTest) addKind(k *g.Kind) (ret *g.Kind) {
	n.ks = append(n.ks, k)
	return k
}

func (n *recordTest) GetKindByName(name string) (ret *g.Kind, err error) {
	var ok bool
	for _, k := range n.ks {
		if k.Name() == name {
			ret, ok = k, true
			break
		}
	}
	if !ok {
		err = errutil.New("kind not found", name)
	}
	return
}

func (n *recordTest) GetField(target, field string) (ret g.Value, err error) {
	switch target {
	case meta.Variables:
		if v, ok := n.vars[field]; !ok {
			err = g.UnknownField(target, field)
		} else {
			ret = g.RecordOf(v)
		}
	default:
		err = errutil.New("unknown field", target, field)
	}
	return
}

type recordTest struct {
	testutil.PanicRuntime
	ks   []*g.Kind
	vars map[string]*g.Record
}

func newRecordAccessTest() *recordTest {
	n := new(recordTest)
	ks := n.addKind(g.NewKindWithTraits(n, "Ks", nil, []g.Field{
		{"d", affine.Number, "" /*"float64"*/},
		{"t", affine.Text, "" /*"string"*/},
		{"a", affine.Text, "a"},
	}, []g.Aspect{{
		Name:   "a",
		Traits: []string{"x", "w", "y"},
	}}))
	n.vars = map[string]*g.Record{
		"boop": ks.NewRecord(),
		"beep": ks.NewRecord(),
	}
	return n
}
