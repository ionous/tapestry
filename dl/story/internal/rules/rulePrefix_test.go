package rules_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story/internal/rules"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
)

func TestPrefixing(t *testing.T) {
	var ks testutil.Kinds
	type Patterns struct{}
	type Misc struct{ Patterns }
	//
	type Actions struct{ Patterns }
	type Events struct{ Patterns }
	//
	type Testing struct{ Actions }
	type BeforeTesting struct{ Events }
	type AfterTesting struct{ Events }

	ks.AddKinds(
		(*Patterns)(nil),
		(*Misc)(nil),
		(*Actions)(nil),
		(*Events)(nil),
		(*Testing)(nil),
		(*BeforeTesting)(nil),
		(*AfterTesting)(nil),
	)

	// test misc:
	// "before" and "after" match misc with appropriate rank
	// "instead of" and "report" dont match.
	t.Run("simple patterns", func(t *testing.T) {
		if e := match(t, &ks, "misc", "misc", 0); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "before misc", "misc", -1); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "after misc", "misc", 1); e != nil {
			t.Fatal(e)
		} else if e := fail(t, &ks, "instead of misc"); e != nil {
			t.Fatal(e)
		} else if e := fail(t, &ks, "report misc"); e != nil {
			t.Fatal(e)
		}
	})

	// test testing:
	// before and after are actual patterns
	// instead and report match before and after with rank
	t.Run("action events", func(t *testing.T) {
		if e := match(t, &ks, "testing", "testing", 0); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "before testing", "before testing", -1); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "after testing", "after testing", 1); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "instead of testing", "before testing", -2); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "report testing", "after testing", 2); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("doubles", func(t *testing.T) {
		if e := failDoubles(t, &ks, "misc"); e != nil {
			t.Fatal(e)
		} else if e := failDoubles(t, &ks, "testing"); e != nil {
			t.Fatal(e)
		}
	})

	type InsteadOfMisc struct{ Patterns }
	type BeforeMisc struct{ Patterns }
	type AfterMisc struct{ Patterns }
	type ReportMisc struct{ Patterns }

	ks.AddKinds(
		(*InsteadOfMisc)(nil),
		(*BeforeMisc)(nil),
		(*AfterMisc)(nil),
		(*ReportMisc)(nil),
	)

	// these behave as they behave; not terribly worried about it.
	t.Run("now what!?", func(t *testing.T) {
		if e := match(t, &ks, "misc", "misc", 0); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "before misc", "misc", -1); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "after misc", "misc", 1); e != nil {
			t.Fatal(e)
		} else if e := fail(t, &ks, "instead of misc"); e != nil {
			t.Fatal(e)
		} else if e := fail(t, &ks, "report misc"); e != nil {
			t.Fatal(e)
		} else if e := failDoubles(t, &ks, "misc"); e == nil {
			t.Fatal("expected success")
		} else {
			t.Log("ok:", e)
		}
	})
}

type logger interface{ Logf(fmt string, args ...any) }

func failDoubles(l logger, ks *testutil.Kinds, test string) (err error) {
	xs := []string{"before", "after", "report"}
	for i := range xs {
		for j := range xs {
			s := strings.Join([]string{xs[i], xs[j], test}, " ")
			if e := fail(l, ks, s); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

func fail(l logger, ks *testutil.Kinds, test string) (err error) {
	if _, e := rules.ReadName(ks, test); e == nil {
		err = errutil.New("expected failure for", test)
	} else {
		l.Logf("ok: %q failed with %s", test, e)
	}
	return
}

func match(l logger, ks *testutil.Kinds, test, name string, rank int) (err error) {
	if p, e := rules.ReadName(ks, test); e != nil {
		err = e
	} else if p.Name != name {
		err = errutil.Fmt("test %q: got name %q wanted %q", test, p.Name, name)
	} else if p.Rank != rank {
		err = errutil.Fmt("test %q: got rank %d wanted %d", test, p.Rank, rank)
	} else {
		l.Logf("ok %q got %q rank %d", test, p.Name, p.Rank)
	}
	return
}
