package raw

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"git.sr.ht/~ionous/tapestry/qna/query"
)

// verify that query none implements every method
var _ query.Query = (*Data)(nil)

// NotImplemented - generic error used returned by QueryNone
type NotImplemented string

func (e NotImplemented) Error() string {
	return string(e)
}

func (q *Data) Close() {
	*q = Data{} // sure, why not.
	return
}

func (q *Data) IsDomainActive(name string) (okay bool, err error) {
	okay = slices.Contains(q.Scenes, name)
	return
}

func (q *Data) ActivateDomains(name string) (prev, next []string, err error) {
	if cnt := len(q.Scenes); cnt == 0 {
		err = errors.New("data has no domains")
	} else if len(name) == 0 {
		// reset
	} else if last := q.Scenes[cnt-1]; name != last {
		err = fmt.Errorf("unexpected scene request %q", name)
	} else {
		next = q.Scenes
	}
	return
}

func (q *Data) ReadChecks(actuallyJustThisOne string) (_ []query.CheckData, _ error) {
	return // none
}

func (q *Data) FieldsOf(exactKind string) (ret []query.FieldData, err error) {
	if i, ok := slices.BinarySearchFunc(q.Kinds, exactKind, func(k KindData, _ string) int {
		return cmp.Compare(k.Kind, exactKind)
	}); !ok {
		err = fmt.Errorf("fields of %q not found", exactKind)
	} else {
		ret = q.Kinds[i].Fields
	}
	return
}

func (q *Data) KindOfAncestors(singleOrPlural string) (ret []string, err error) {
	if kind, e := q.getPluralKind(singleOrPlural); e != nil {
		err = e
	} else if i, ok := slices.BinarySearchFunc(q.Kinds, kind, func(k KindData, _ string) int {
		return cmp.Compare(k.Kind, kind)
	}); !ok {
		err = fmt.Errorf("fields of %q not found", kind)
	} else {
		ret = q.Kinds[i].Ancestors
	}
	return
}

// search using short name
func (q *Data) NounInfo(short string) (ret query.NounInfo, err error) {
	if noun, e := q.findbyShort(short); e != nil {
		err = e
	} else if n, e := q.findByFull(noun); e != nil {
		err = e
	} else {
		ret = query.NounInfo{
			Domain: n.Domain,
			Noun:   n.Noun,
			Kind:   n.Kind,
		}
	}
	return
}

// the best short name for the passed full name
func (q *Data) NounName(full string) (ret string, err error) {
	if noun, e := q.findByFull(full); e != nil {
		err = e
	} else if as := noun.Aliases; len(as) == 0 {
		ret = full
	} else {
		ret = as[0]
	}
	return
}

func (q *Data) NounNames(full string) (ret []string, err error) {
	if noun, e := q.findByFull(full); e != nil {
		err = e
	} else {
		ret = noun.Aliases
	}
	return
}

func (q *Data) NounValues(full, field string) (ret []query.ValueData, err error) {
	if noun, e := q.findByFull(full); e != nil {
		err = e
	} else {
		start, end := -1, len(noun.Values)
		for i, nv := range noun.Values {
			at := nv.Field == field
			if at && start < 0 {
				start = i
			} else if !at && start >= 0 {
				end = i
				break
			}
		}
		if start >= 0 {
			ret = noun.Values[start:end]
		}
	}
	return
}

func (q *Data) NounsByKind(kind string) (ret []string, _ error) {
	for _, n := range q.Nouns { // tbd: record nouns per kind?
		if n.Kind == kind {
			ret = append(ret, n.Noun)
		}
	}
	return
}

func (q *Data) PluralToSingular(plural string) (ret string, _ error) {
	for _, p := range q.Plurals {
		if plural == p.Many {
			ret = p.One
			break
		}
	}
	return
}

func (q *Data) PluralFromSingular(singular string) (ret string, _ error) {
	if i, ok := slices.BinarySearchFunc(q.Plurals, singular, func(n Plural, _ string) int {
		return cmp.Compare(n.One, singular)
	}); ok {
		ret = q.Plurals[i].Many
	}
	return
}

func (q *Data) PatternLabels(pat string) (ret []string, err error) {
	if p, e := q.findPattern(pat); e != nil {
		err = e
	} else {
		ret = p.Labels
	}
	return
}

func (q *Data) RulesFor(pat string) (ret []query.RuleData, err error) {
	if p, e := q.findPattern(pat); e != nil {
		err = e
	} else {
		ret = p.Rules
	}
	return
}

func (q *Data) ReciprocalsOf(rel, noun string) (ret []string, err error) {
	if p, e := q.findRelation(rel); e != nil {
		err = e
	} else {
		for _, el := range p.Pairs {
			if el.Other == noun {
				ret = append(ret, el.One)
			}
		}
	}
	return
}

func (q *Data) RelativesOf(rel, noun string) (ret []string, err error) {
	if p, e := q.findRelation(rel); e != nil {
		err = e
	} else {
		for _, el := range p.Pairs {
			if el.One == noun {
				ret = append(ret, el.Other)
			}
		}
	}
	return
}

// FIX!
func (q *Data) Relate(rel, noun, otherNoun string) error {
	// if p, e := q.findRelation(rel); e != nil {
	// 	err = e
	// } else {
	// 	for _, el := range p.Pairs {
	// 		if el.One == noun {
	// 			ret = append(ret, el.Other)
	// 		}
	// 	}
	// }
	// return
	panic("not implemented")
}

// Random implements Query.
func (q *Data) Random(inclusiveMin int, exclusiveMax int) int {
	// FIX!!!!
	return inclusiveMin
}

// LoadGame implements Query.
func (q *Data) LoadGame(path string) (ret query.CacheMap, err error) {
	err = NotImplemented("load game")
	return
}

// SaveGame implements Query.
func (q *Data) SaveGame(path string, dynamicValues query.CacheMap) error {
	return NotImplemented("save game")
}

// normalize plural name
func (q *Data) getPluralKind(singleOrPlural string) (ret string, err error) {
	if n, e := q.PluralFromSingular(singleOrPlural); e != nil {
		err = e
	} else if len(n) == 0 {
		ret = singleOrPlural
	} else {
		ret = n
	}
	return
}

// shortname to id
func (q *Data) findbyShort(shortname string) (ret string, err error) {
	if i, ok := slices.BinarySearchFunc(q.Names, shortname, func(n NounName, _ string) int {
		return cmp.Compare(n.Name, shortname)
	}); !ok {
		err = fmt.Errorf("couldnt find noun with shortname %q", shortname)
	} else {
		ret = q.Names[i].Noun
	}
	return
}

// fullname to NounData
func (q *Data) findByFull(fullname string) (ret NounData, err error) {
	if i, ok := slices.BinarySearchFunc(q.Nouns, fullname, func(n NounData, _ string) int {
		return cmp.Compare(n.Noun, fullname)
	}); !ok {
		err = fmt.Errorf("couldnt find noun with fullname %q", fullname)
	} else {
		ret = q.Nouns[i]
	}
	return
}

func (q *Data) findPattern(name string) (ret PatternData, err error) {
	if i, ok := slices.BinarySearchFunc(q.Patterns, name, func(p PatternData, _ string) int {
		return cmp.Compare(p.Pattern, name)
	}); !ok {
		err = fmt.Errorf("couldnt find pattern %q", name)
	} else {
		ret = q.Patterns[i]
	}
	return
}

func (q *Data) findRelation(name string) (ret RelativeData, err error) {
	if i, ok := slices.BinarySearchFunc(q.Relatives, name, func(p RelativeData, _ string) int {
		return cmp.Compare(p.Relation, name)
	}); !ok {
		err = fmt.Errorf("couldnt find relation %q", name)
	} else {
		ret = q.Relatives[i]
	}
	return
}
