package story

import (
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/ephemera/eph"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/tables"
)

func NewActionName(k *Importer, n ActionName) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_PATTERN, n.At.String()), nil
}

func NewAspect(k *Importer, n Aspect) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_ASPECT, n.At.String()), nil
}

func NewEventName(k *Importer, n EventName) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_EVENT, n.At.String()), nil
}

func NewNounName(k *Importer, n NounName) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_NOUN, n.At.String()), nil
}

func NewNameWithCategory(k *Importer, n NounName, cat string) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, cat, n.At.String()), nil
}

func NewPatternName(k *Importer, n value.PatternName) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_PATTERN, n.At.String()), nil
}

func NewPluralKinds(k *Importer, n PluralKinds) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_PLURAL_KINDS, n.At.String()), nil
}

func FixSingular(k *Importer, n PluralKinds) (ret eph.Named, err error) {
	name := lang.Breakcase(lang.Singularize(n.Str))
	return k.NewName(name, tables.NAMED_KIND, n.At.String()), nil
}

func NewProperty(k *Importer, n Property) (ret eph.Named, err error) {
	// note: this is linked to NAMED_ASPECT
	// aspect properties in kinds currently must have the same name as the aspect.
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_FIELD, n.At.String()), nil
}

func NewRecordSingular(k *Importer, n RecordSingular) (ret eph.Named, err error) {
	// fix? for now, we leverage the existing kind assembly
	// name := lang.LowerBreakcase(n.Str)
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_KIND, n.At.String()), nil
}

func NewRecordPlural(k *Importer, n RecordPlural) (ret eph.Named, err error) {
	// fix? for now, we leverage the existing kind assembly
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_PLURAL_KINDS, n.At.String()), nil
}

func NewRelation(k *Importer, n value.RelationName) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_VERB, n.At.String()), nil
}

func NewSingularKind(k *Importer, n SingularKind) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_KIND, n.At.String()), nil
}

// fix: this is *not* where pluralization should happen
func FixPlurals(k *Importer, n SingularKind) (ret eph.Named, err error) {
	name := lang.Breakcase(lang.Pluralize(n.Str))
	return k.NewName(name, tables.NAMED_KIND, n.At.String()), nil
}

func NewTestName(k *Importer, n TestName) (ret eph.Named, err error) {
	// fix? all names should probably munge their own strings
	// ( see ephemera's NewDomainName for the current hack )
	// things that would need work are:
	// tests, autogen fields, and control over the domain ( ie. NewName vs. NewDowmainName )
	if n.Str == TestName_CurrentTest {
		ret = k.Env().Recent.Test
	} else {
		name := lang.Breakcase(n.Str)
		ret = k.NewName(name, tables.NAMED_TEST, n.At.String())
	}
	return
}

func NewTrait(k *Importer, n Trait) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, tables.NAMED_TRAIT, n.At.String()), nil
}

func NewVariableName(k *Importer, n value.VariableName, cat string) (ret eph.Named, err error) {
	name := lang.Breakcase(n.Str)
	return k.NewName(name, cat, n.At.String()), nil
}
