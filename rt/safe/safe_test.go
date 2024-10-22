package safe_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestSafety(t *testing.T) {
	var run testutil.PanicRuntime
	switch e := safe.RunAll(&run, nil); e.(type) {
	case nil:
		t.Log("okay nothing run")
	default:
		t.Fatal(e)
	}
	switch e := safe.Run(&run, nil); e.(type) {
	case safe.MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
	switch _, e := safe.GetBool(&run, nil); e.(type) {
	case safe.MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
	switch _, e := safe.GetNum(&run, nil); e.(type) {
	case safe.MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
	switch _, e := safe.GetText(&run, nil); e.(type) {
	case safe.MissingEval:
		t.Log("okay", e)
	default:
		t.Fatal(e)
	}
}
