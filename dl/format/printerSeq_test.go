package format

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestSequences(t *testing.T) {
	t.Run("cycle none", func(t *testing.T) {
		matchSequence(t, []string{
			"",
		}, &CycleText{
			Name: t.Name(),
		})
	})
	t.Run("cycle text", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "b", "c", "a", "b", "c", "a",
		}, &CycleText{
			Name: t.Name(), Parts: []rt.TextEval{
				literal.T("a"),
				literal.T("b"),
				literal.T("c"),
			}})
	})
	t.Run("stopping", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "b", "c", "c", "c", "c", "c",
		}, &StoppingText{
			Name: t.Name(), Parts: []rt.TextEval{
				literal.T("a"),
				literal.T("b"),
				literal.T("c"),
			}})
	})
	t.Run("once", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "", "", "", "",
		}, &StoppingText{
			Name: t.Name(), Parts: []rt.TextEval{
				literal.T("a"),
			}})
	})
	t.Run("shuffle one", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "a",
		}, &ShuffleText{
			Name: t.Name(), Parts: []rt.TextEval{
				literal.T("a"),
			}})
	})
	t.Run("shuffle", func(t *testing.T) {
		matchSequence(t, []string{
			"c", "d", "b", "e", "a", "b", "e",
		}, &ShuffleText{
			Name: t.Name(), Parts: []rt.TextEval{
				literal.T("a"),
				literal.T("b"),
				literal.T("c"),
				literal.T("d"),
				literal.T("e"),
			}})
	})
}

func matchSequence(t *testing.T, want []string, seq rt.TextEval) {
	run := seqTest{counters: make(map[string]int)}
	var have []string
	for i, wanted := range want {
		if got, e := safe.GetText(&run, seq); e != nil {
			t.Fatal(e)
		} else if got := got.String(); got != wanted {
			t.Fatalf("error at %d wanted %q got %q", i, wanted, got)
		} else {
			have = append(have, got)
		}
	}
	t.Log(t.Name(), have)
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type seqTest struct {
	baseRuntime
	counters map[string]int
}

func (m *seqTest) Random(inclusiveMin, exclusiveMax int) int {
	return (exclusiveMax-inclusiveMin)/2 + inclusiveMin
}

func (m *seqTest) GetField(target, field string) (ret rt.Value, err error) {
	if target != meta.Counter {
		err = rt.UnknownField(target, field)
	} else {
		v := m.counters[field]
		ret = rt.IntOf(v)
	}
	return
}

func (m *seqTest) SetField(target, field string, value rt.Value) (err error) {
	if target != meta.Counter {
		err = rt.UnknownField(target, field)
	} else {
		m.counters[field] = value.Int()
	}
	return
}
