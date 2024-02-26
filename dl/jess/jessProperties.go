package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func (op *KindsHaveProperties) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kind.Match(q, &next) &&
		op.Have.Match(q, &next, keywords.Have, keywords.Has) &&
		(Optional(q, &next, &op.Article) || true) &&
		(op.matchListOf(q, &next) || true) &&
		op.PropertyType.Match(q, &next) {
		Optional(q, &next, &op.CalledName)
		*input, okay = next, true
	}
	return
}

func (op *KindsHaveProperties) matchListOf(q Query, input *InputState) (okay bool) {
	if m, width := listOf.FindMatch(input.Words()); m != nil {
		op.ListOf = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

// register a single field
func (op *KindsHaveProperties) Generate(rar Registrar) (err error) {
	if kind, e := op.Kind.Validate(kindsOf.Kind, kindsOf.Record); e != nil {
		err = e
	} else if f, e := op.PropertyType.GetType(op.ListOf != nil); e != nil {
		err = e
	} else {
		if op.CalledName != nil {
			f.Name = op.CalledName.String()
		}
		if len(f.Name) == 0 {
			// erroring feels like more useful than failing to match...
			err = fmt.Errorf("%s fields require an explicit name", f.Affinity)
		} else {
			err = rar.AddFields(kind, []mdl.FieldInfo{f})
		}
	}
	return
}

func (op *PropertyType) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.matchPrimitive(q, &next) ||
		Optional(q, &next, &op.Kind) {
		*input, okay = next, true
	}
	return
}

func (op *PropertyType) matchPrimitive(q Query, input *InputState) (okay bool) {
	if m, width := primitiveTypes.FindMatch(input.Words()); m != nil {
		op.Primitive = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

// return a default field name, its affine type and its optional class name
func (op *PropertyType) GetType(listOf bool) (ret mdl.FieldInfo, err error) {

	var name string
	var aff affine.Affinity
	var cls string

	if p := op.Primitive; p != nil {
		aff, cls = getTypeOfPrim(p.String())
	} else {
		// use the name the author specified for the field
		name = inflect.Normalize(op.Kind.Matched.String())
		// even if that differs from the actual name of the kind...
		aff, cls, err = getTypeOfKind(op.Kind)
	}
	if err == nil {
		// bump up the affinity to list
		// tbd? could list affinity be a flag instead?
		// (well... affinity would have to be a set of const first)
		if listOf {
			switch aff {
			case affine.Number:
				aff = affine.NumList
			case affine.Text:
				aff = affine.TextList
			case affine.Record:
				aff = affine.RecordList
			}
		}
		ret = mdl.FieldInfo{Name: name, Affinity: aff, Class: cls}
	}
	return
}

// in jess, class name always empty for primitive types
func getTypeOfPrim(str string) (retAff affine.Affinity, retCls string) {
	// fix? inform supports "texts" or "text"
	// jess only supports "text", and it cheats
	// assuming if not text, then number/s.
	if str == affine.Text.String() {
		retAff = affine.Text
	} else {
		retAff = affine.Number
	}
	return
}

// - kind and aspect will generate affine text
// - record will generate affine record
func getTypeOfKind(k *Kind) (retAff affine.Affinity, retCls string, err error) {
	kt := k.ActualKind.base
	switch kt {
	case kindsOf.Kind, kindsOf.Aspect:
		retAff = affine.Text
		retCls = k.String()
	case kindsOf.Record:
		retAff = affine.Record
		retCls = k.String()
	default:
		err = fmt.Errorf("unexpected kind of property %q", kt)
	}
	return
}

var listOf = match.PanicSpans("list of")
var primitiveTypes = match.PanicSpans("text", "numbers", "number")
