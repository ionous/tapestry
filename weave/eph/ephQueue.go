package eph

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

type EphemeraWriter interface{ WriteEphemera(Ephemera) }

type WriterFun func(eph Ephemera)

func (w WriterFun) WriteEphemera(op Ephemera) {
	w(op)
}

func NewCommandBuilder(Queue EphemeraWriter) *CommandBuilder {
	return &CommandBuilder{Queue: Queue}
}

func NewCommandQueue(els *[]Ephemera) *CommandBuilder {
	var w WriterFun = func(Queue Ephemera) {
		(*els) = append(*els, Queue)
	}
	return &CommandBuilder{Queue: w}
}

type CommandBuilder struct {
	Queue   EphemeraWriter
	domains []string
}

func (k *CommandBuilder) append(p Ephemera) {
	k.Queue.WriteEphemera(p)
}

func (k *CommandBuilder) BeginDomain(name string, requires []string) (none error) {
	k.domains = append(k.domains, name)
	k.append(&BeginDomain{
		Name:     name,
		Requires: requires,
	})
	return
}

func (k *CommandBuilder) EndDomain() (err error) {
	if top := len(k.domains) - 1; top < 0 {
		err = errutil.New("unexpected end domain")
	} else {
		name := k.domains[top]
		k.domains = k.domains[:top]
		k.append(&EndDomain{
			Name: name,
		})
	}
	return
}

func (k *CommandBuilder) AssertAlias(short string, names ...string) (none error) {
	k.append(&Aliases{
		ShortName: short,
		Aliases:   names,
	})
	return
}

func (k *CommandBuilder) AssertAncestor(kind, ancestor string) (none error) {
	if ancestor == kindsOf.Pattern.String() {
		// fix: this should be possible to replace.
		k.append(&Patterns{PatternName: kind})
	} else {
		k.append(&Kinds{
			Kind:     kind,
			Ancestor: ancestor,
		})
	}
	return
}

func (k *CommandBuilder) AssertAspectTraits(aspect string, traits []string) (none error) {
	k.append(&Aspects{
		Aspects: aspect,
		Traits:  traits,
	})
	return
}

func (k *CommandBuilder) AssertCheck(name string, do []rt.Execute, expect literal.LiteralValue) (none error) {
	k.append(&Checks{
		Name:   name,
		Exe:    do,
		Expect: expect,
	})
	return
}

// make an arbitrary string key and value
// the key includes all elements of the path except for the final one
// which gets used as the value.
// differing definitions generation conflicts.
func (k *CommandBuilder) AssertDefinition(path ...string) (err error) {
	if end := len(path) - 1; end < 1 {
		err = errutil.New("missing key value pair for definition", path)
	} else {
		k.append(&Definition{
			Path:  path[:end],
			Value: path[end],
		})
	}
	return
}

func makeParam(name, class string, aff affine.Affinity, init assign.Assignment) (ret Params, err error) {
	if init != nil {
		if test := assign.GetAffinity(init); test != aff {
			err = errutil.Fmt(`mismatched affinity of initial value (a %s) for field "%s" (a %s)`, test, name, aff)
		}
	}
	if err == nil {
		ret = Params{
			Name:      name,
			Affinity:  affineToAffinity(aff),
			Class:     class,
			Initially: init,
		}
	}
	return
}

func (k *CommandBuilder) AssertField(kind, name, class string, aff affine.Affinity, init assign.Assignment) (err error) {
	if ps, e := makeParam(name, class, aff, init); e != nil {
		err = e
	} else {
		k.append(&Kinds{
			Kind:    kind,
			Contain: []Params{ps},
		})
	}
	return
}

func (k *CommandBuilder) AssertParam(kind, name, class string, aff affine.Affinity, init assign.Assignment) (err error) {
	if ps, e := makeParam(name, class, aff, init); e != nil {
		err = e
	} else {
		k.append(&Patterns{
			PatternName: kind,
			Params:      []Params{ps},
		})
	}
	return
}

func (k *CommandBuilder) AssertLocal(kind, name, class string, aff affine.Affinity, init assign.Assignment) (err error) {
	if ps, e := makeParam(name, class, aff, init); e != nil {
		err = e
	} else {
		k.append(&Patterns{
			PatternName: kind,
			Locals:      []Params{ps},
		})
	}
	return
}

func (k *CommandBuilder) AssertResult(kind, name, class string, aff affine.Affinity, init assign.Assignment) (err error) {
	if ps, e := makeParam(name, class, aff, init); e != nil {
		err = e
	} else {
		k.append(&Patterns{
			PatternName: kind,
			Result:      &ps,
		})
	}
	return
}

func (k *CommandBuilder) AssertGrammar(name string, d *grammar.Directive) (none error) {
	k.append(&Directives{
		Name:      name,
		Directive: *d,
	})
	return
}

func (k *CommandBuilder) AssertNounKind(noun, kind string) (none error) {
	k.append(&Nouns{
		Noun: noun,
		Kind: kind,
	})
	return
}

// func (k *CommandBuilder) AssertNounPhrase() (none error) {
// 	panic("not implemented")
// 	return
// }

func (k *CommandBuilder) AssertNounValue(noun, field string, path []string, val literal.LiteralValue) (none error) {
	k.append(&Values{
		Noun:  noun,
		Field: field,
		Path:  path,
		Value: val,
	})
	return
}

func (k *CommandBuilder) AssertOpposite(a, b string) (none error) {
	k.append(&Opposites{
		Opposite: a,
		Word:     b,
	})
	return
}

func (k *CommandBuilder) AssertPlural(singluar, plural string) (none error) {
	k.append(&Plurals{
		Singular: singluar,
		Plural:   plural,
	})
	return
}

func (k *CommandBuilder) AssertRelation(rel, a, b string, amany, bmany bool) (err error) {
	switch {
	case amany && bmany:
		err = k.assertRelation(rel,
			&ManyMany{Kinds: a, OtherKinds: b},
		)
	case !amany && !bmany:
		err = k.assertRelation(rel,
			&OneOne{Kind: a, OtherKind: b},
		)
	case amany && !bmany:
		err = k.assertRelation(rel,
			&ManyOne{Kinds: a, OtherKind: b},
		)
	case !amany && bmany:
		err = k.assertRelation(rel,
			&OneMany{Kind: a, OtherKinds: b},
		)
	default:
		panic("stray neutrino detected")
	}
	return
}

func (k *CommandBuilder) assertRelation(rel string, card Cardinality) (none error) {
	k.append(&Relations{
		Rel:         rel,
		Cardinality: card,
	})
	return
}

func (k *CommandBuilder) AssertRelative(rel, noun, otherNoun string) (none error) {
	k.append(&Relatives{
		Rel:       rel,
		Noun:      noun,
		OtherNoun: otherNoun,
	})
	return
}
func (k *CommandBuilder) AssertRule(name string, target string, guard rt.BoolEval, flags assert.EventTiming, do []rt.Execute) (none error) {
	t, a := fromTiming(flags)
	k.append(&Rules{
		PatternName: name,
		Target:      target,
		Filter:      guard,
		When:        t,
		Exe:         do,
		Touch:       a,
	})
	return
}

func (k *CommandBuilder) Schedule(when assert.Phase, do func(assert.World, assert.Assertions) error) (none error) {
	k.append(PhaseFunction{when, func(nm assert.World, nk assert.Assertions) error {
		return do(nm, nk)
	}})
	return
}

// translate "bool" to "$BOOL", etc.
// note: can return affine.None ( unknown affinity )
func affineToAffinity(a affine.Affinity) (ret Affinity) {
	spec := ret.Compose()
	if k, i := spec.IndexOfValue(a.String()); i >= 0 {
		ret.Str = k
	}
	return
}
