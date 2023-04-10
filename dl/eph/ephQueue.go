package eph

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

type EphemeraWriter interface{ WriteEphemera(Ephemera) }

type Queue struct {
	q       EphemeraWriter
	domains []string
}

func (k *Queue) append(p Ephemera) {
	k.q.WriteEphemera(p)
}

func (k *Queue) BeginDomain(name string, requires []string) (none error) {
	k.domains = append(k.domains, name)
	k.append(&EphBeginDomain{
		Name:     name,
		Requires: requires,
	})
	return
}

func (k *Queue) EndDomain() (err error) {
	if top := len(k.domains) - 1; top < 0 {
		err = errutil.New("unexpected end domain")
	} else {
		name := k.domains[top]
		k.domains = k.domains[:top]
		k.append(&EphEndDomain{
			Name: name,
		})
	}
	return
}

func (k *Queue) AssertAlias(short string, names ...string) (none error) {
	k.append(&EphAliases{
		ShortName: short,
		Aliases:   names,
	})
	return
}

func (k *Queue) AssertAncestor(kind, ancestor string) (none error) {
	k.append(&EphKinds{
		Kind:     kind,
		Ancestor: ancestor,
	})
	return
}

func (k *Queue) AssertAspectTraits(aspect string, traits []string) (none error) {
	k.append(&EphAspects{
		Aspects: aspect,
		Traits:  traits,
	})
	return
}

func (k *Queue) AssertCheck(name string, do []rt.Execute) (none error) {
	k.append(&EphChecks{
		Name: name,
		Exe:  do,
	})
	return
}

func (k *Queue) AssertDefinition(path ...string) (none error) {
	panic("not implemented")
	// k.append(&EphDefinition{
	// 	Path: path,
	// })
	return
}

func makeParam(name, class string, init assign.Assignment) EphParams {
	return EphParams{
		Name:      name,
		Affinity:  affineToAffinity(assign.GetAffinity(init)),
		Class:     class,
		Initially: init,
	}
}

func (k *Queue) AssertField(kind, name, class string, init assign.Assignment) (err error) {
	ps := makeParam(name, class, init)
	k.append(&EphKinds{
		Kind:    kind,
		Contain: []EphParams{ps},
	})
	return
}

func (k *Queue) AssertParam(kind, name, class string, init assign.Assignment) (err error) {
	ps := makeParam(name, class, init)
	k.append(&EphPatterns{
		PatternName: kind,
		Params:      []EphParams{ps},
	})
	return
}

func (k *Queue) AssertLocal(kind, name, class string, init assign.Assignment) (err error) {
	ps := makeParam(name, class, init)
	k.append(&EphPatterns{
		PatternName: kind,
		Locals:      []EphParams{ps},
	})
	return
}

func (k *Queue) AssertResult(kind, name, class string, init assign.Assignment) (err error) {
	ps := makeParam(name, class, init)
	k.append(&EphPatterns{
		PatternName: kind,
		Result:      &ps,
	})
	return
}

func (k *Queue) AssertGrammar(name string, d *grammar.Directive) (none error) {
	k.append(&EphDirectives{
		Name:      name,
		Directive: *d,
	})
	return
}

// func (k *Catalog) AssertMacro( a, b string ) {
//    // 40: 		k.WriteEphemera(&EphMacro{EphPatterns: out, MacroStatements: op.MacroStatements})
//  }

func (k *Queue) AssertNounKind(noun, kind string) (none error) {
	k.append(&EphNouns{
		Noun: noun,
		Kind: kind,
	})
	return
}

func (k *Queue) AssertNounPhrase() (none error) {
	panic("not implemented")
	return
}

func (k *Queue) AssertNounValue(noun, field string, path []string, val literal.LiteralValue) (none error) {
	k.append(&EphValues{
		Noun:  noun,
		Field: field,
		Path:  path,
		Value: val,
	})
	return
}

func (k *Queue) AssertOpposite(a, b string) (none error) {
	k.append(&EphOpposites{
		Opposite: a,
		Word:     b,
	})
	return
}

func (k *Queue) AssertPlural(singluar, plural string) (none error) {
	k.append(&EphPlurals{
		Singular: singluar,
		Plural:   plural,
	})
	return
}

func (k *Queue) AssertRelation(rel, a, b string, amany, bmany bool) (err error) {
	switch {
	case amany && bmany:
		err = k.assertRelation(rel, EphCardinality{
			EphCardinality_ManyMany_Opt,
			&ManyMany{Kinds: a, OtherKinds: b},
		})
	case !amany && !bmany:
		err = k.assertRelation(rel, EphCardinality{
			EphCardinality_OneOne_Opt,
			&OneOne{Kind: a, OtherKind: b},
		})
	case amany && !bmany:
		err = k.assertRelation(rel, EphCardinality{
			EphCardinality_ManyOne_Opt,
			&ManyOne{Kinds: a, OtherKind: b},
		})
	case !amany && bmany:
		err = k.assertRelation(rel, EphCardinality{
			EphCardinality_OneMany_Opt,
			&OneMany{Kind: a, OtherKinds: b},
		})
	default:
		panic("stray neutrino detected")
	}
	return
}

func (k *Queue) assertRelation(rel string, card EphCardinality) (none error) {
	k.append(&EphRelations{
		Rel:         rel,
		Cardinality: card,
	})
	return
}

func (k *Queue) AssertRelative(rel, noun, otherNoun string) (none error) {
	k.append(&EphRelatives{
		Rel:       rel,
		Noun:      noun,
		OtherNoun: otherNoun,
	})
	return
}
func (k *Queue) AssertRule(name string, target string, guard rt.BoolEval, flags assert.EventTiming, do []rt.Execute) (none error) {
	t, a := fromTiming(flags)
	k.append(&EphRules{
		PatternName: name,
		Target:      target,
		Filter:      guard,
		When:        t,
		Exe:         do,
		Touch:       a,
	})
	return
}

func (k *Queue) Schedule(when assert.Phase, do func() error) (none error) {
	k.append(PhaseFunction{when, func(*Catalog, *Domain, string) error {
		return do()
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
