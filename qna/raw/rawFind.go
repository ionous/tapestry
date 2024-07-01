package raw

import (
	"cmp"
	"slices"

	"git.sr.ht/~ionous/tapestry/rt"
)

func FindKind(ks []rt.Kind, exactKind string) (ret *rt.Kind, okay bool) {
	if i, ok := slices.BinarySearchFunc(ks, exactKind, func(k rt.Kind, _ string) int {
		return cmp.Compare(k.Name(), exactKind)
	}); ok {
		ret, okay = &ks[i], true
	}
	return
}
func FindRelation(els []RelativeData, name string) (ret *RelativeData, okay bool) {
	if i, ok := slices.BinarySearchFunc(els, name, func(k RelativeData, _ string) int {
		return cmp.Compare(k.Relation, name)
	}); ok {
		ret, okay = &(els[i]), true
	}
	return
}
func FindValueField(els []EvalData, name string) (ret EvalData, okay bool) {
	if i, ok := slices.BinarySearchFunc(els, name, func(k EvalData, _ string) int {
		return cmp.Compare(k.Field, name)
	}); ok {
		ret, okay = els[i], true
	}
	return
}
func FindRecordField(els []RecordData, name string) (ret RecordData, okay bool) {
	if i, ok := slices.BinarySearchFunc(els, name, func(k RecordData, _ string) int {
		return cmp.Compare(k.Field, name)
	}); ok {
		ret, okay = els[i], true
	}
	return
}
func FindPlural(els []Plural, name string) (ret Plural, okay bool) {
	if i, ok := slices.BinarySearchFunc(els, name, func(k Plural, _ string) int {
		return cmp.Compare(k.One, name)
	}); ok {
		ret, okay = els[i], true
	}
	return
}
func FindName(els []NounName, name string) (ret NounName, okay bool) {
	if i, ok := slices.BinarySearchFunc(els, name, func(k NounName, _ string) int {
		return cmp.Compare(k.Name, name)
	}); ok {
		ret, okay = els[i], true
	}
	return
}
func FindNoun(els []NounData, name string) (ret NounData, okay bool) {
	if i, ok := slices.BinarySearchFunc(els, name, func(k NounData, _ string) int {
		return cmp.Compare(k.Noun, name)
	}); ok {
		ret, okay = els[i], true
	}
	return
}
func FindPattern(els []PatternData, name string) (ret PatternData, okay bool) {
	if i, ok := slices.BinarySearchFunc(els, name, func(k PatternData, _ string) int {
		return cmp.Compare(k.Pattern, name)
	}); ok {
		ret, okay = els[i], true
	}
	return
}
