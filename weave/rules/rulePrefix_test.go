package rules_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/test/testutil"
	"git.sr.ht/~ionous/tapestry/weave/rules"
)

func TestPrefixing(t *testing.T) {
	var ks testutil.Kinds
	type Patterns struct{}
	type Misc struct{ Patterns }
	//
	type Actions struct{ Patterns }
	//
	type Testing struct{ Actions }
	type BeforeTesting struct{ Testing }
	type AfterTesting struct{ Testing }

	ks.AddKinds(
		(*Patterns)(nil),
		(*Misc)(nil),
		(*Actions)(nil),
		(*Testing)(nil),
		(*BeforeTesting)(nil),
		(*AfterTesting)(nil),
	)

	var (
		before  = rules.Ranks[0]
		instead = rules.Ranks[1]
		at      = rules.Ranks[2]
		after   = rules.Ranks[3]
		report  = rules.Ranks[4]
	)

	// test misc:
	// "before" and "after" match misc with appropriate rank
	// "instead of" and "report" dont match.
	t.Run("simple patterns", func(t *testing.T) {
		if e := match(t, &ks, "misc", "misc", at); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "before misc", "misc", before); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "after misc", "misc", after); e != nil {
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
		} else if e := match(t, &ks, "before testing", "before testing", before); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "after testing", "after testing", after); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "instead of testing", "before testing", instead); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "report testing", "after testing", report); e != nil {
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
		} else if e := match(t, &ks, "before misc", "misc", before); e != nil {
			t.Fatal(e)
		} else if e := match(t, &ks, "after misc", "misc", after); e != nil {
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

func failDoubles(l logger, ks *testutil.Kinds, phrase string) (err error) {
	var errs []error
	xs := []string{"before", "after", "report"}
	for i := range xs {
		for j := range xs {
			s := strings.Join([]string{xs[i], xs[j], phrase}, " ")
			if e := fail(l, ks, s); e != nil {
				errs = append(errs, e)
			}
		}
	}
	return errors.Join(errs...)
}

func fail(l logger, ks *testutil.Kinds, phrase string) (err error) {
	if n, e := rules.ReadPhrase(ks, phrase, ""); e != nil {
		err = e
	} else if _, e := n.GetRuleInfo(); e != nil {
		err = e
	}
	//
	if err != nil {
		l.Logf("ok: %q failed with %s", phrase, err)
		err = nil
	} else {
		err = fmt.Errorf("expected %s failure for get rule info", phrase)
	}
	return
}

// /*else

func match(l logger, ks *testutil.Kinds, phrase, name string, rank int) (err error) {
	if n, e := rules.ReadPhrase(ks, phrase, ""); e != nil {
		err = e
	} else if p, e := n.GetRuleInfo(); e != nil {
		err = e
	} else if p.Name != name {
		err = fmt.Errorf("test %q: got name %q wanted %q", phrase, p.Name, name)
	} else if p.Rank != rank {
		err = fmt.Errorf("test %q: got rank %d wanted %d", phrase, p.Rank, rank)
	} else {
		l.Logf("ok %q got %q rank %d", phrase, p.Name, p.Rank)
	}
	return
}
