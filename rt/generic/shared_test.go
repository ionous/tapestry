package generic_test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
)

func (n *recordTest) addKind(k *rt.Kind) (ret *rt.Kind) {
	n.ks = append(n.ks, k)
	return k
}

func (n *recordTest) GetKindByName(name string) (ret *rt.Kind, err error) {
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

func (n *recordTest) GetField(target, field string) (ret rt.Value, err error) {
	switch target {
	case meta.Variables:
		if v, ok := n.vars[field]; !ok {
			err = rt.UnknownField(target, field)
		} else {
			ret = rt.RecordOf(v)
		}
	default:
		err = errutil.New("unknown field", target, field)
	}
	return
}

type recordTest struct {
	testutil.PanicRuntime
	ks   []*rt.Kind
	vars map[string]*rt.Record
}

func newRecordAccessTest() *recordTest {
	n := new(recordTest)
	ks := n.addKind(rt.NewKind([]string{"Ks"}, []rt.Field{
		{Name: "d", Affinity: affine.Number},
		{Name: "t", Affinity: affine.Text},
		{Name: "a", Affinity: affine.Text, Type: "a"},
	}, []rt.Aspect{{
		Name:   "a",
		Traits: []string{"x", "w", "y"},
	}}))
	n.vars = map[string]*rt.Record{
		"boop": rt.NewRecord(ks),
		"beep": rt.NewRecord(ks),
	}
	return n
}
