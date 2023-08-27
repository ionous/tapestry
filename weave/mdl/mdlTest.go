package mdl

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"github.com/ionous/errutil"
)

var LogWarning = func(e error) {
	log.Println("Warning:", e) // for now good enough
}

type Warnings []error

// override the global warning function
// returns a defer-able function which:
// 1. restores the warning function; and,
// 2. raises a Fatal error if there are any unhandled warnings.
func (w *Warnings) Catch(fatal func(args ...any)) func() {
	was := LogWarning
	LogWarning = func(e error) {
		(*w) = append((*w), e)
	}
	return func() {
		if len(*w) > 0 {
			fatal("unhandled warnings", *w)
		}
		LogWarning = was
	}
}

// return the warnings as a raw list, clear all stored errors.
func (w *Warnings) All() (ret []error) {
	ret, (*w) = (*w), nil
	return ret
}

// remove and return the first warning,
// error if there is none, or if warning doesnt start with the passed prefix.
func (w *Warnings) Expect(prefix string) (err error) {
	if e := w.pop(); e == nil {
		err = errutil.Fmt("expected %q, received nothing.", prefix)
	} else if str := e.Error(); !strings.HasPrefix(str, prefix) {
		err = errutil.Fmt("expected %q, received %q.", prefix, str)
	}
	return
}

func (w *Warnings) pop() (err error) {
	if cnt := len(*w); cnt > 0 {
		err, (*w) = (*w)[0], (*w)[1:]
	}
	return
}

// for testing: a generic field of the kind
func (pen *Pen) AddTestField(kind, field string, aff affine.Affinity, cls string) (err error) {
	if kid, e := pen.findRequiredKind(kind); e != nil {
		err = errutil.Fmt("%w trying to add field %q", e, field)
	} else if cls, e := pen.findOptionalKind(cls); e != nil {
		err = errutil.Fmt("%w trying to write field %q", e, field)
	} else {
		e := pen.addField(kid, cls, field, aff)
		err = eatDuplicates(pen.warn, e)
	}
	return
}

func (pen *Pen) AddTestParameter(kind, field string, aff affine.Affinity, cls string) (err error) {
	if kid, e := pen.findRequiredKind(kind); e != nil {
		err = errutil.Fmt("%w trying to add parameter %q", e, field)
	} else if cls, e := pen.findOptionalKind(cls); e != nil {
		err = errutil.Fmt("%w trying to write parameter %q", e, field)
	} else {
		err = pen.addParameter(kid, cls, field, aff)
	}
	return
}

func (pen *Pen) AddTestResult(kind, field string, aff affine.Affinity, cls string) (err error) {
	if kid, e := pen.findRequiredKind(kind); e != nil {
		err = errutil.Fmt("%w trying to add parameter %q", e, field)
	} else if cls, e := pen.findOptionalKind(cls); e != nil {
		err = errutil.Fmt("%w trying to write parameter %q", e, field)
	} else {
		err = pen.addResult(kid, cls, field, aff)
	}
	return
}

// public for tests:
func (pen *Pen) AddTestRule(pattern string, rank int, prog string) (err error) {
	domain, at := pen.domain, pen.at
	if kid, e := pen.findRequiredKind(pattern); e != nil {
		err = e
	} else {
		var name any = nil
		_, err = pen.db.Exec(mdl_rule, domain, kid.id, name, rank, 0, 0, 0, prog, at)
	}
	return
}

// unmarshaled version of AddValue for testing.
func (pen *Pen) AddTestValue(noun, path, out string) (err error) {
	if noun, e := pen.findRequiredNoun(noun, nounWithKind); e != nil {
		err = e
	} else {
		parts := strings.Split(path, ".")
		if outer, _, e := pen.digField(noun, parts); e != nil {
			err = e // for testing, we accept any inner most affinity ( so long as the parts were resolvable )
		} else {
			root, dot := parts[0], strings.Join(parts[1:], ".")
			err = pen.addValue(noun, outer, root, dot, out)
		}
	}
	return
}
