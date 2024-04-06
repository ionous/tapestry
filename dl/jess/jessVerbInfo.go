package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// a snapshot of the verb; represented as a noun in the db.
type VerbDesc struct {
	Subject,
	Object,
	Alternate,
	Relation string
	Implies  []string
	Reversed bool
}

func (v *VerbDesc) applyVerb(u Scheduler, w weaver.Weaves, lhs, rhs []DesiredNoun) (err error) {
	// do some extra work to always generate the nouns on the left hand side first
	lk, rk := v.getKinds()
	if e := addKindToNouns(w, lk, lhs); e != nil {
		err = e
	} else if e := addKindToNouns(w, rk, rhs); e != nil {
		err = e
	} else {
		// then pick the dependent side for implications and pairs
		var subjects, objects []DesiredNoun
		if !v.Reversed {
			subjects, objects = lhs, rhs
		} else {
			subjects, objects = rhs, lhs
		}
		if ts := v.Implies; len(ts) > 0 {
			for i := range objects {
				objects[i].appendTraits(ts)
			}
		}
		if e := writePairs(w, v.Relation, subjects, objects); e != nil {
			err = e
		} else if len(v.Alternate) > 0 {
			// if there was an alternate kind for the subjects
			// then we only set the noun to "Objects"
			// and we need to fallback to something more specific
			// at least one of the two kinds *must* successfully be applied.
			err = generateFallbacks(u, subjects, v.Subject, v.Alternate)
		}
	}
	return
}

func (v *VerbDesc) getKinds() (lhs, rhs string) {
	// when "alternate" is set -- only mark the object as an object
	// we'll do a pass to ensure all is well.
	var subject string
	if len(v.Alternate) > 0 {
		subject = Objects
	} else {
		subject = v.Subject
	}
	if !v.Reversed {
		lhs, rhs = subject, v.Object
	} else {
		rhs, lhs = subject, v.Object
	}
	return
}

func addKindToNouns(w weaver.Weaves, kind string, ns []DesiredNoun) (err error) {
	for _, n := range ns {
		if e := w.AddNounKind(n.Noun, kind); e != nil && !errors.Is(e, weaver.Duplicate) {
			err = e
			break
		}
	}
	return
}

func writePairs(w weaver.Weaves, rel string, ps, cs []DesiredNoun) (err error) {
Pairs:
	for _, p := range ps {
		for _, c := range cs {
			if e := w.AddNounPair(rel, p.Noun, c.Noun); e != nil {
				err = e
				break Pairs
			}
		}
	}
	return
}

func readVerb(run rt.Runtime, verb string) (ret VerbDesc, err error) {
	if relation, e := readString(run, verb, VerbRelation); e != nil {
		err = e
	} else if object, e := readString(run, verb, VerbObject); e != nil {
		err = e
	} else if subject, e := readString(run, verb, VerbSubject); e != nil {
		err = e
	} else if alternate, e := readString(run, verb, VerbAlternate); e != nil && !errors.Is(e, weaver.Missing) {
		err = e // alternate subects(s) are optional
	} else if implies, e := readStringList(run, verb, VerbImplies); e != nil && !errors.Is(e, weaver.Missing) {
		err = e // implications are optional
	} else if rev, revErr := readString(run, verb, VerbReversed); revErr != nil && !errors.Is(revErr, weaver.Missing) {
		err = revErr // reverse is optional; false if not explicitly specified
	} else {
		ret = VerbDesc{
			Subject:   subject,
			Object:    object,
			Alternate: alternate,
			Relation:  relation,
			// these are normally specified as raw text
			// any errors show up when they are applied
			Implies:  implies,
			Reversed: rev == ReversedTrait,
		}
	}
	return
}

func readString(run rt.Runtime, noun, field string) (ret string, err error) {
	if b, e := run.GetField(noun, field); e != nil {
		err = e
	} else if aff := b.Affinity(); aff != affine.Text {
		err = fmt.Errorf(`expected that "%s.%s" was text, not %s`, noun, field, aff)
	} else {
		ret = b.String()
	}
	return
}

func readStringList(run rt.Runtime, noun, field string) (ret []string, err error) {
	if b, e := run.GetField(noun, field); e != nil {
		err = e
	} else if aff := b.Affinity(); aff != affine.TextList {
		err = fmt.Errorf(`expected that "%s.%s" was a list of text, not %s`, noun, field, aff)
	} else {
		ret = b.Strings()
	}
	return
}
