package jess

import (
	"encoding/json"
	"errors"

	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type verbInfo struct {
	subject,
	object,
	alternate,
	relation,
	implication string
	reversed bool
}

func (v *verbInfo) applyVerb(ctx *Context, lhs, rhs []DesiredNoun) (err error) {
	// do some extra work to always generate the nouns on the left hand side first
	lk, rk := v.getKinds()
	if e := addKindToNouns(ctx, lk, lhs); e != nil {
		err = e
	} else if e := addKindToNouns(ctx, rk, rhs); e != nil {
		err = e
	} else {
		// then pick the dependent side for implications and pairs
		var subjects, objects []DesiredNoun
		if !v.reversed {
			subjects, objects = lhs, rhs
		} else {
			subjects, objects = rhs, lhs
		}
		if trait := v.implication; len(trait) > 0 {
			for i := range objects {
				objects[i].appendTrait(trait)
			}
		}
		if e := writePairs(ctx, v.relation, subjects, objects); e != nil {
			err = e
		} else if len(v.alternate) > 0 {
			// if there was an alternate kind for the subjects
			// then we only set the noun to "Objects"
			// and we need to fallback to something more specific
			// at least one of the two kinds *must* successfully be applied.
			err = generateFallbacks(ctx, subjects, v.subject, v.alternate)
		}
	}
	return
}

func (v *verbInfo) getKinds() (lhs, rhs string) {
	// when "alternate" is set -- only mark the object as an object
	// we'll do a pass to ensure all is well.
	var subject string
	if len(v.alternate) > 0 {
		subject = Objects
	} else {
		subject = v.subject
	}
	if !v.reversed {
		lhs, rhs = subject, v.object
	} else {
		rhs, lhs = subject, v.object
	}
	return
}

func addKindToNouns(ctx *Context, kind string, ns []DesiredNoun) (err error) {
	for _, n := range ns {
		if e := ctx.AddNounKind(n.Noun, kind); e != nil && !errors.Is(e, mdl.Duplicate) {
			err = e
			break
		}
	}
	return
}

func writePairs(ctx *Context, rel string, ps, cs []DesiredNoun) (err error) {
Pairs:
	for _, p := range ps {
		for _, c := range cs {
			if e := ctx.AddNounPair(rel, p.Noun, c.Noun); e != nil {
				err = e
				break Pairs
			}
		}
	}
	return
}

func readVerb(ctx *Context, verb string) (ret verbInfo, err error) {
	if relation, e := ctx.readString(verb, VerbRelation); e != nil {
		err = e
	} else if object, e := ctx.readString(verb, VerbObject); e != nil {
		err = e
	} else if subject, e := ctx.readString(verb, VerbSubject); e != nil {
		err = e
	} else if alternate, e := ctx.readString(verb, VerbAlternate); e != nil && !errors.Is(e, mdl.Missing) {
		err = e // alternate subects(s) are optional
	} else if implication, e := ctx.readString(verb, VerbImplication); e != nil && !errors.Is(e, mdl.Missing) {
		err = e // implication(s) are optional
	} else if rev, revErr := ctx.readString(verb, VerbReversed); revErr != nil && !errors.Is(revErr, mdl.Missing) {
		err = revErr // reverse is optional; false if not explicitly specified
	} else {
		ret = verbInfo{
			subject:     subject,
			object:      object,
			alternate:   alternate,
			relation:    relation,
			implication: implication,
			reversed:    rev == ReversedTrait,
		}
	}
	return
}

func (ctx *Context) readString(noun, field string) (ret string, err error) {
	if b, e := ctx.GetNounValue(noun, field); e != nil {
		err = e
	} else {
		err = json.Unmarshal(b, &ret)
	}
	return

}
