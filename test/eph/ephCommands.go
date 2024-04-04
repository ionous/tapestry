package eph

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type Aliases struct {
	ShortName string
	Aliases   []string
}

func (op *Aliases) Assert(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if n, e := run.GetField(meta.ObjectId, op.ShortName); e != nil {
			err = e
		} else {
			n := n.String()
			for _, a := range op.Aliases {
				if e := w.AddNounName(n, a, -1); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

// Aspects A set of related object states such that exactly one member of the set is true for a given object at a single time.
// Generates an implicit kind of 'aspect' where every field of the kind is a boolean property.
type Aspects struct {
	Aspects string
	Traits  []string
}

func (op *Aspects) Assert(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.AncestryPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if e := w.AddKind(op.Aspects, kindsOf.Aspect.String()); e != nil {
			err = e
		} else {
			// no reason not to add the traits immediately since we have them already
			err = w.AddAspectTraits(op.Aspects, op.Traits)
		}
		return
	})
}

type BeginDomain struct {
	Name     string
	Requires []string
}

func (op *BeginDomain) Assert(cat *weave.Catalog) (err error) {
	return cat.DomainStart(op.Name, op.Requires)
}

// Directives
type Directives struct {
	Directive grammar.Directive
}

func (op *Directives) Assert(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
		return w.AddGrammar(op.Directive.Name, &op.Directive)
	})
}

// EndDomain
type EndDomain struct {
	Name string
}

func (op *EndDomain) Assert(cat *weave.Catalog) (err error) {
	return cat.DomainEnd()
}

// Kinds A new type deriving from another existing type.
// The new kind has all of the properties of all of its ancestor kinds
// and it can be used wherever one of its ancestor kinds is needed.
// ( The reverse isn't true because the new kind can have its own unique properties not available to its ancestors. )
type Kinds struct {
	Kind     string
	Ancestor string
	Contain  []Params
}

func (op *Kinds) Assert(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.AncestryPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if e := w.AddKind(op.Kind, op.Ancestor); e != nil {
			err = e
		} else {
			err = cat.Schedule(weaver.PropertyPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
				if ps := op.Contain; len(ps) > 0 {
					err = w.AddKindFields(op.Kind, reduceFields(ps))
				}
				return
			})
		}
		return
	})
}

// Nouns
type Nouns struct {
	Noun string
	Kind string
}

func (op *Nouns) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weaver.NounPhase, func(w weaver.Weaves, run rt.Runtime) error {
		_, err := mdl.AddNamedNoun(w, op.Noun, op.Kind)
		return err
	})
}

// Patterns Patterns provide author reusable code.
// The parameters define values provided by the caller.
// Locals provide scratch values for use during pattern processing.
// The result allows the pattern to return a value to the caller of pattern.
// While multiple pattern commands can be used to define a pattern,
// the set of arguments and the return can only be specified once.
type Patterns struct {
	PatternName string
	Params      []Params
	Locals      []Params
	Result      *Params
}

func (op *Patterns) Assert(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		kb := mdl.NewPatternBuilder(op.PatternName)
		kb.AddLocals(reduceFields(op.Locals))
		kb.AddParams(reduceFields(op.Params))
		if p := op.Result; p != nil {
			kb.AddResult(p.GetFieldInfo())
		}
		return w.AddPattern(kb.Pattern)
	})
}

// Plurals Rules for transforming plural text to singular text and back again.
// Used by the assembler to help interpret author definitions,
// and at runtime to help the parser interpret user input.
type Plurals struct {
	Plural   string
	Singular string
}

func (op *Plurals) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weaver.LanguagePhase, func(w weaver.Weaves, run rt.Runtime) error {
		return w.AddPlural(op.Plural, op.Singular)
	})
}

// Relations
type Relations struct {
	Rel         string
	Cardinality Cardinality
}

func (op *Relations) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weaver.PropertyPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		switch c := op.Cardinality.(type) {
		case *OneOne:
			err = w.AddRelation(op.Rel, c.Kind, c.OtherKind, false, false)
		case *OneMany:
			err = w.AddRelation(op.Rel, c.Kind, c.OtherKinds, false, true)
		case *ManyOne:
			err = w.AddRelation(op.Rel, c.Kinds, c.OtherKind, true, false)
		case *ManyMany:
			err = w.AddRelation(op.Rel, c.Kinds, c.OtherKinds, true, true)
		}
		return
	})
}

// Relatives
type Relatives struct {
	Rel       string
	Noun      string
	OtherNoun string
}

func (op *Relatives) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weaver.ConnectionPhase, func(w weaver.Weaves, run rt.Runtime) error {
		return w.AddNounPair(op.Rel, op.Noun, op.OtherNoun)
	})
}

// Rules
type Rules struct {
	PatternName string
	Exe         []rt.Execute
}

func (op *Rules) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
		pb := mdl.NewPatternBuilder(op.PatternName)
		pb.AppendRule(0, rt.Rule{
			Exe: op.Exe,
		})
		return w.ExtendPattern(pb.Pattern)
	})
}

// Values Give a noun a specific value at startup.
// Initialization is somewhat simplistic:
// 1. Initial values are not scoped to domains, triggers must be used to change values when domains begin and end.
// 2. The values inside of records can be set using a 'path' to find them, however individual values within lists cannot be set.
// Note: when using a path, the path addresses the noun first, the named field - referring to the inner most record - last.
type Values struct {
	Noun  string
	Field string
	Path  []string
	Value literal.LiteralValue
}

func (op *Values) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if n, e := run.GetField(meta.ObjectId, op.Noun); e != nil {
			err = e
		} else {
			n := n.String()
			if field, path := op.Field, op.Path; len(path) == 0 {
				err = w.AddNounValue(n, field, assign.Literal(op.Value))
			} else {
				path := append(path, field)
				err = w.AddNounPath(n, path, op.Value)
			}
		}
		return
	})
}
