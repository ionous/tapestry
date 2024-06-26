package raw

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"

	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
)

// verify that query none implements every method
var _ query.Query = (*RawQuery)(nil)

// NotImplemented - generic error used returned by QueryNone
type NotImplemented string

func (e NotImplemented) Error() string {
	return string(e)
}

// fix: still needs the decoder for record values
// but it could be minimized to the literal decoder
func MakeQuery(data *Data, dec decoder.Decoder) RawQuery {
	return RawQuery{data, dec}
}

type RawQuery struct {
	*Data
	dec decoder.Decoder
}

func (q RawQuery) Close() {
	*q.Data = Data{} // sure, why not.
	return
}

func (q RawQuery) IsDomainActive(name string) (okay bool, err error) {
	okay = slices.Contains(q.Scenes, name)
	return
}

func (q RawQuery) ActivateDomains(name string) (prev, next []string, err error) {
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

func (q RawQuery) ReadChecks(actuallyJustThisOne string) (_ []query.CheckData, _ error) {
	return // none
}

func (q RawQuery) KindOfAncestors(singleOrPlural string) (ret []string, err error) {
	if exactKind, e := q.getPluralKind(singleOrPlural); e != nil {
		err = e
	} else if k, e := q.GetKindByName(exactKind); e != nil {
		err = e
	} else {
		ret = k.Ancestors()
	}
	return
}

// search using short name
func (q RawQuery) NounInfo(short string) (ret query.NounInfo, err error) {
	if noun, e := q.findFullName(short); e != nil {
		err = e
	} else if n, e := q.findNounData(noun); e != nil {
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
func (q RawQuery) NounName(full string) (ret string, err error) {
	if noun, e := q.findNounData(full); e != nil {
		err = e
	} else {
		if name := noun.CommonName; len(name) > 0 {
			ret = name
		} else {
			ret = noun.Noun
		}
	}
	return
}

func (q RawQuery) NounNames(full string) (ret []string, err error) {
	if noun, e := q.findNounData(full); e != nil {
		err = e
	} else {
		ret = noun.Aliases
	}
	return
}

func (q RawQuery) NounValue(full, field string) (ret rt.Assignment, err error) {
	panic("fix")
	// if noun, e := q.findNounData(full); e != nil {
	// 	err = e
	// } else {
	// 	start, end := -1, len(noun.Values)
	// 	for i, nv := range noun.Values {
	// 		at := nv.Field == field
	// 		if at && start < 0 {
	// 			start = i
	// 		} else if !at && start >= 0 {
	// 			end = i
	// 			break
	// 		}
	// 	}
	// 	if start >= 0 {
	// 		ret = noun.Values[start:end]
	// 	}
	// }
	// return
}

func (q RawQuery) NounsByKind(kind string) (ret []string, _ error) {
	for _, n := range q.Nouns { // tbd: record nouns per kind?
		if n.Kind == kind {
			ret = append(ret, n.Noun)
		}
	}
	return
}

func (q RawQuery) PluralToSingular(plural string) (ret string, _ error) {
	for _, p := range q.Plurals {
		if plural == p.Many {
			ret = p.One
			break
		}
	}
	return
}

func (q RawQuery) PluralFromSingular(singular string) (ret string, _ error) {
	if i, ok := slices.BinarySearchFunc(q.Plurals, singular, func(n Plural, _ string) int {
		return cmp.Compare(n.One, singular)
	}); ok {
		ret = q.Plurals[i].Many
	}
	return
}

func (q RawQuery) PatternLabels(pat string) (ret []string, err error) {
	if p, e := q.findPattern(pat); e != nil {
		err = e
	} else {
		ret = p.Labels
	}
	return
}

func (q RawQuery) RulesFor(pat string) (ret query.RuleSet, err error) {
	if p, e := q.findPattern(pat); e != nil {
		err = e
	} else {
		ret = query.RuleSet{
			Rules:     p.Rules,
			UpdateAll: p.UpdateAll,
		}
	}
	return
}

func (q RawQuery) ReciprocalsOf(rel, noun string) (ret []string, err error) {
	if p, e := q.findRelation(rel); e != nil {
		err = e
	} else {
		for _, el := range p.Pairs {
			if el.Other == noun && len(el.One) > 0 {
				ret = append(ret, el.One)
			}
		}
	}
	return
}

func (q RawQuery) RelativesOf(rel, noun string) (ret []string, err error) {
	if p, e := q.findRelation(rel); e != nil {
		err = e
	} else {
		for _, el := range p.Pairs {
			if el.One == noun && len(el.Other) > 0 {
				ret = append(ret, el.Other)
			}
		}
	}
	return
}

func (q RawQuery) Relate(rel, noun, otherNoun string) (err error) {
	if p, e := q.findRelation(rel); e != nil {
		err = e
	} else if one, other := len(noun) > 0, len(otherNoun) > 0; !one && !other {
		err = errors.New("nothing to relate")
	} else {
		// follows from the sqlite version: relatePair
		suffix := strings.HasSuffix(p.Cardinality, "one")
		prefix := strings.HasPrefix(p.Cardinality, "one")
		if (!suffix && !prefix) && (!one || !other) {
			err = errors.New("not implemented: erasing from many-many relationships needs more thought")
		} else {
			els := slices.DeleteFunc(p.Pairs, func(prev Pair) bool {
				return (suffix && prev.One == noun) ||
					(prefix && prev.Other == otherNoun)
			})
			// if both sides exist: we're adding a new pair
			if one && other {
				els = append(els, Pair{noun, otherNoun})
			}
			// find relation returns a pointer
			p.Pairs = els
		}
	}
	return
}

// Random implements Query.
func (q RawQuery) Random(inclusiveMin int, exclusiveMax int) int {
	// FIX!!!!
	return inclusiveMin
}

// LoadGame implements Query.
func (q RawQuery) LoadGame(path string) (ret query.CacheMap, err error) {
	err = NotImplemented("load game")
	return
}

// SaveGame implements Query.
func (q RawQuery) SaveGame(path string, dynamicValues query.CacheMap) error {
	return NotImplemented("save game")
}

// normalize plural name
func (q RawQuery) getPluralKind(singleOrPlural string) (ret string, err error) {
	if n, e := q.PluralFromSingular(singleOrPlural); e != nil {
		err = e
	} else if len(n) == 0 {
		ret = singleOrPlural
	} else {
		ret = n
	}
	return
}

func (q *Data) GetKindByName(exactKind string) (*rt.Kind, error) {
	return FindKind(q.Kinds, exactKind)
}

// shortname to id
func (q RawQuery) findFullName(shortname string) (ret string, err error) {
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
func (q RawQuery) findNounData(fullname string) (ret NounData, err error) {
	if i, ok := slices.BinarySearchFunc(q.Nouns, fullname, func(n NounData, _ string) int {
		return cmp.Compare(n.Noun, fullname)
	}); !ok {
		err = fmt.Errorf("couldnt find noun with fullname %q", fullname)
	} else {
		ret = q.Nouns[i]
	}
	return
}

func (q RawQuery) findPattern(name string) (ret PatternData, err error) {
	if i, ok := slices.BinarySearchFunc(q.Patterns, name, func(p PatternData, _ string) int {
		return cmp.Compare(p.Pattern, name)
	}); !ok {
		err = fmt.Errorf("couldnt find pattern %q", name)
	} else {
		ret = q.Patterns[i]
	}
	return
}

func (q RawQuery) findRelation(name string) (ret *RelativeData, err error) {
	if i, ok := slices.BinarySearchFunc(q.Relatives, name, func(p RelativeData, _ string) int {
		return cmp.Compare(p.Relation, name)
	}); !ok {
		err = fmt.Errorf("couldnt find relation %q", name)
	} else {
		ret = &(q.Relatives[i])
	}
	return
}
