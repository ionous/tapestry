package eph

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type Aliases struct {
	ShortName string
	Aliases   []string
}

func (op *Aliases) Assert(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequireAll, func(w *weave.Weaver) (err error) {
		n := lang.Normalize(op.ShortName)
		if n, e := w.GetClosestNoun(n); e != nil {
			err = e
		} else {
			pen := w.Pin()
			for _, a := range op.Aliases {
				if e := pen.AddName(n, a, -1); e != nil {
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
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) (err error) {
		return w.Pin().AddAspect(op.Aspects, op.Traits)
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
	return cat.Schedule(weave.RequireRules, func(w *weave.Weaver) error {
		return w.Pin().AddGrammar(op.Directive.Name, &op.Directive)
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
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) (err error) {
		pen := w.Pin()
		if e := pen.AddKind(op.Kind, op.Ancestor); e != nil {
			err = e
		} else if ps := op.Contain; len(ps) > 0 {
			fields := mdl.NewFieldBuilder(op.Kind)
			for _, p := range ps {
				fields.AddField(p.FieldInfo())
			}
			err = pen.AddFields(fields.Fields)
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
	return cat.Schedule(weave.RequireDefaults, func(w *weave.Weaver) error {
		return w.Pin().AddNoun(op.Noun, op.Noun, op.Kind)
	})
}

// Opposites Rules for transforming plural text to singular text and back again.
// Used by the assembler to help interpret author definitions,
// and at runtime to help the parser interpret user input.
type Opposites struct {
	Opposite string
	Word     string
}

func (op *Opposites) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) error {
		return w.Pin().AddOpposite(op.Opposite, op.Word)
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
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) (err error) {
		kb := mdl.NewPatternBuilder(op.PatternName)
		if ps := op.Params; err == nil && len(ps) > 0 {
			for _, p := range ps {
				kb.AddParam(p.FieldInfo())
			}
		}
		if p := op.Result; err == nil && p != nil {
			kb.AddResult(p.FieldInfo())
		}
		if ps := op.Locals; err == nil {
			for _, p := range ps {
				kb.AddLocal(p.FieldInfo())
			}
		}
		return w.Pin().AddPattern(kb.Pattern)
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
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) error {
		return w.Pin().AddPlural(op.Plural, op.Singular)
	})
}

// Relations
type Relations struct {
	Rel         string
	Cardinality Cardinality
}

func (op *Relations) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) (err error) {
		pen := w.Pin()
		switch c := op.Cardinality.(type) {
		case *OneOne:
			err = pen.AddRelation(op.Rel, c.Kind, c.OtherKind, false, false)
		case *OneMany:
			err = pen.AddRelation(op.Rel, c.Kind, c.OtherKinds, false, true)
		case *ManyOne:
			err = pen.AddRelation(op.Rel, c.Kinds, c.OtherKind, true, false)
		case *ManyMany:
			err = pen.AddRelation(op.Rel, c.Kinds, c.OtherKinds, true, true)
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
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) error {
		return w.Pin().AddPair(op.Rel, op.Noun, op.OtherNoun)
	})
}

// Rules
type Rules struct {
	PatternName string
	Exe         []rt.Execute
}

func (op *Rules) Assert(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) error {
		pb := mdl.NewPatternBuilder(op.PatternName)
		pb.AppendRule(0, rt.Rule{
			Exe: op.Exe,
		})
		return w.Pin().ExtendPattern(pb.Pattern)
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
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) (err error) {
		pen := w.Pin()
		if n, e := pen.GetClosestNoun(op.Noun); e != nil {
			err = e
		} else if field, path := op.Field, op.Path; len(path) == 0 {
			err = pen.AddInitialValue(n, field, assign.Literal(op.Value))
		} else {
			path := append(path, field)
			err = pen.AddPathValue(n, mdl.MakePath(path...), op.Value)
		}
		return
	})
}
