package core

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func TestSequences(t *testing.T) {
	t.Run("cycle none", func(t *testing.T) {
		matchSequence(t, []string{
			"",
		}, &CallCycle{
			Name: t.Name(),
		})
	})
	t.Run("cycle text", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "b", "c", "a", "b", "c", "a",
		}, &CallCycle{
			Name: t.Name(), Parts: []rt.TextEval{
				T("a"),
				T("b"),
				T("c"),
			}})
	})
	t.Run("stopping", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "b", "c", "c", "c", "c", "c",
		}, &CallTerminal{
			Name: t.Name(), Parts: []rt.TextEval{
				T("a"),
				T("b"),
				T("c"),
			}})
	})
	t.Run("once", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "", "", "", "",
		}, &CallTerminal{
			Name: t.Name(), Parts: []rt.TextEval{
				T("a"),
			}})
	})
	t.Run("shuffle one", func(t *testing.T) {
		matchSequence(t, []string{
			"a", "a",
		}, &CallShuffle{
			Name: t.Name(), Parts: []rt.TextEval{
				T("a"),
			}})
	})
	t.Run("shuffle", func(t *testing.T) {
		matchSequence(t, []string{
			"c", "d", "b", "e", "a", "b", "e",
		}, &CallShuffle{
			Name: t.Name(), Parts: []rt.TextEval{
				T("a"),
				T("b"),
				T("c"),
				T("d"),
				T("e"),
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

type seqTest struct {
	baseRuntime
	counters map[string]int
}

func (m *seqTest) Random(inclusiveMin, exclusiveMax int) int {
	return (exclusiveMax-inclusiveMin)/2 + inclusiveMin
}

func (m *seqTest) GetField(target, field string) (ret g.Value, err error) {
	if target != meta.Counter {
		err = g.UnknownField(target, field)
	} else {
		v := m.counters[field]
		ret = g.IntOf(v)
	}
	return
}

func (m *seqTest) SetField(target, field string, value g.Value) (err error) {
	if target != meta.Counter {
		err = g.UnknownField(target, field)
	} else {
		m.counters[field] = value.Int()
	}
	return
}
