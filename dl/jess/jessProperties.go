package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *CalledName) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Called.Match(q, &next) &&
		op.Name.Match(AddContext(q, CheckIndefiniteArticles), &next) {
		// fix: does this need CheckIndefiniteArticles here?
		// that seems a little weird, because when does this generate a noun
		// to have that applied?
		*input, okay = next, true
	}
	return
}

func (op *CalledName) GetNormalizedName() (string, error) {
	return op.Name.GetNormalizedName()
}

// runs in the PropertyPhase
func (op *KindsHaveProperties) Phase() weaver.Phase {
	return weaver.PropertyPhase
}

func (op *KindsHaveProperties) MatchLine(q Query, line InputState) (ret InputState, okay bool) {
	if next := line; //
	op.Kind.Match(q, &next) &&
		op.Have.Match(q, &next, keywords.Have, keywords.Has) &&
		(Optional(q, &next, &op.Article) || true) &&
		(op.matchListOf(&next) || true) &&
		op.PropertyType.Match(q, &next) {
		Optional(q, &next, &op.CalledName)
		ret, okay = next, true
	}
	return
}

func (op *KindsHaveProperties) matchListOf(input *InputState) (okay bool) {
	if m, width := listOf.FindPrefix(input.Words()); m != nil {
		op.ListOf = m.String()
		*input = input.Skip(width)
		okay = true
	}
	return
}

// register a single field
func (op *KindsHaveProperties) Weave(w weaver.Weaves, run rt.Runtime) (err error) {
	if kind, e := op.Kind.Validate(kindsOf.Kind, kindsOf.Record); e != nil {
		err = e
	} else if f, e := op.PropertyType.GetType(len(op.ListOf) > 0); e != nil {
		err = e
	} else {
		// has the author explicitly overridden the default field name?
		if op.CalledName != nil {
			f.Name, err = op.CalledName.GetNormalizedName()
		}
		if err == nil {
			if len(f.Name) == 0 {
				// erroring feels like more useful than failing to match...
				err = fmt.Errorf("%s fields require an explicit name", f.Affinity)
			} else {
				err = w.AddKindFields(kind, []mdl.FieldInfo{f})
			}
		}
	}
	return
}

func (op *PropertyType) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.matchPrimitive(&next) ||
		Optional(q, &next, &op.Kind) {
		*input, okay = next, true
	}
	return
}

func (op *PropertyType) matchPrimitive(input *InputState) (okay bool) {
	if m, width := primitiveTypes.FindPrefix(input.Words()); m != nil {
		op.Primitive = m.String()
		*input = input.Skip(width)
		okay = true
	}
	return
}

// return a default field name, its affine type and its optional class name
func (op *PropertyType) GetType(listOf bool) (ret mdl.FieldInfo, err error) {
	var name, cls string
	var aff affine.Affinity
	if prim := op.Primitive; len(prim) > 0 {
		aff, cls = getTypeOfPrim(prim)
	} else if n, e := match.NormalizeAll(op.Kind.Matched); e != nil {
		err = e
	} else {
		name = n
		// re: singular or plural; use the author specified field name
		// even if that differs from the actual name of the kind...
		aff, cls, err = getTypeOfKind(op.Kind)
	}
	if err == nil {
		// bump up the affinity to list
		// tbd? could list affinity be a flag instead?
		// (well... affinity would have to be a set of const first)
		if listOf {
			switch aff {
			case affine.Num:
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
		retAff = affine.Num
	}
	return
}

// - kind and aspect will generate affine text
// - record will generate affine record
func getTypeOfKind(k *Kind) (retAff affine.Affinity, retCls string, err error) {
	kt := k.actualKind.BaseKind
	switch kt {
	case kindsOf.Kind, kindsOf.Aspect:
		retAff = affine.Text
		retCls = k.actualKind.Name
	case kindsOf.Record:
		retAff = affine.Record
		retCls = k.actualKind.Name
	default:
		err = fmt.Errorf("unexpected kind of property %q", kt)
	}
	return
}

var listOf = match.PanicSpans("list of")
var primitiveTypes = match.PanicSpans("text", "numbers", "number")
