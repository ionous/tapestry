package jess

import (
	"log"

	"git.sr.ht/~ionous/tapestry/support/grok"
)

func (op *Article) Match(q Query, in InputState) (ret int) {
	op.Matched, ret = q.FindArticle(in.Words())
	return
}

func (op *TraitName) Match(q Query, in InputState) (ret int) {
	if trait, width := q.FindTrait(in.Words()); width > 0 {
		op.Matched, ret = trait, width
	}
	return
}

func (op *KindName) Match(q Query, in InputState) (ret int) {
	if trait, width := q.FindKind(in.Words()); width > 0 {
		op.Matched, ret = trait, width
	}
	return
}

func (op *CommaAnd) Match(q Query, in InputState) (ret int) {
	if sep, e := grok.CommaAnd(in.Words()); e != nil {
		log.Println(e) // fix?!
		ret = -1
	} else {
		ret = sep.Len()
	}
	return
}

func (op *Traits) GetTraits() []Matched {
	var out []Matched
	for t := *op; ; {
		out = append(out, t.TraitName.Matched)
		if next := t.AdditionalTraits; next == nil {
			break
		} else {
			t = next.Traits
		}
	}
	return out
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *Traits) Match(q Query, in InputState) (ret int) {
	var article int
	if article = Optionally(q, in, &op.Article); article > 0 {
		in = in.Skip(article)
	}
	if name := op.TraitName.Match(q, in); name > 0 {
		rest := Optionally(q, in.Skip(name), &op.AdditionalTraits)
		ret = article + name + rest
	}
	return
}

func (op *AdditionalTraits) Match(q Query, in InputState) (ret int) {
	var sep int
	if sep = Optionally(q, in, &op.CommaAnd); sep > 0 {
		in = in.Skip(sep)
	}
	if m := op.Traits.Match(q, in); m > 0 {
		ret = sep + m
	}
	return
}

// a generic method to read an optional element.
// optional elements are implemented by pointer types
// the pointers are required to implement Interpreter
func Optionally[M any,
	IM interface {
		// the syntax for this feels very strange
		// the method takes an 'out' pointer ( so go can determine the type by inference )
		// the "interface" is a reused keyword, signifying a type constraint
		// indicating we want pointers to M to also be (implement) Interpreter.
		// *phew*
		*M
		Interpreter
	}](q Query, in InputState, out **M) (ret int) {
	var v M // oh yeah, and then some casting...
	if m := IM(&v).Match(q, in); m > 0 {
		*out, ret = &v, m
	}
	return
}
