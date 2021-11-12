package jsnutil_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/testdl"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/jsnutil"
	"github.com/kr/pretty"
)

func TestT(t *testing.T) {
	flow := testdl.TestFlow{
		Slots: []testdl.TestSlot{
			&testdl.TestFlow{
				Slot: &testdl.TestFlow{
					Txt: testdl.TestTxt{"rock"},
				},
			},
			&testdl.TestFlow{
				Swap: testdl.TestSwap{
					Choice: testdl.TestSwap_B_Opt,
					Value:  &testdl.TestFlow{},
				},
			}},
	}
	var l logger
	if e := jsnutil.Visit(&flow, &l); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(l.out, []string{
		"-> test_flow",
		"-> test_flow/[0]/test_slot/test_flow",
		"-> test_flow/[0]/test_slot/test_flow/test_slot/test_flow",
		"<- test_flow/[0]/test_slot/test_flow/test_slot/test_flow",
		"<- test_flow/[0]/test_slot/test_flow",
		"-> test_flow/[1]/test_slot/test_flow",
		"-> test_flow/[1]/test_slot/test_flow/$B/test_flow",
		"<- test_flow/[1]/test_slot/test_flow/$B/test_flow",
		"<- test_flow/[1]/test_slot/test_flow",
		"<- test_flow",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(l.out))
		t.Fatal(diff)
	}
}

type logger struct {
	out []string
}

func (l *logger) BlockStart(b jsn.FlowBlock, c jsnutil.Context) (ret error) {
	l.out = append(l.out, "-> "+c.Path)
	return
}

func (l *logger) BlockEnd(b jsn.FlowBlock, c jsnutil.Context) (ret error) {
	l.out = append(l.out, "<- "+c.Path)
	return
}
