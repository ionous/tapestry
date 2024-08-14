package raw

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/pack"
)

// verify that raw implements every method of query
var _ query.Query = (*RawQuery)(nil)

func MakeQuery(data *Data) RawQuery {
	rand := query.RandomizedTime()
	return RawQuery{data, &rand}
}

type RawQuery struct {
	*Data
	rand *query.Randomizer
}

func (q RawQuery) Close() {
	*q.Data = Data{} // sure, why not.
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
	if noun, e := q.findNounData(full); e != nil {
		err = e
	} else if a, ok := FindValueField(noun.Values, field); ok {
		ret = a.Value
	} else {
		// otherwise, get the field data
		if k, e := q.GetKindByName(noun.Kind); e != nil {
			err = e
		} else if i := k.FieldIndex(field); i < 0 {
			err = fmt.Errorf("couldnt find field %q in kind %q", field, k.Name())
		} else {
			ft := k.Field(i)
			if a, ok := FindRecordField(noun.Records, field); !ok {
				// if no record: then default value
				ret, err = zeroAssignment(ft, i)
			} else {
				// otherwise, unpack the record and return it
				if rec, e := pack.UnpackRecord(q, a.Packed, ft.Type); e != nil {
					err = e
				} else {
					ret = rt.AssignValue(rt.RecordOf(rec))
				}
			}
		}
	}
	return
}

// ex. give me all nouns of type "actors"
func (q RawQuery) NounsWithAncestor(kind string) (ret []string, err error) {
	for _, n := range q.Nouns {
		// quick check:
		if n.Kind == kind {
			ret = append(ret, n.Noun)
		} else {
			// slow check:
			if nk, e := q.GetKindByName(n.Kind); e != nil {
				err = e
				break
			} else {
				nks := nk.Ancestors() // path: animals, *actors*, things....
				if i := slices.Index(nks, kind); i >= 0 {
					ret = append(ret, n.Noun)
				}
			}
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
	if a, ok := FindPlural(q.Plurals, singular); ok {
		ret = a.Many
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
	return q.rand.Random(inclusiveMin, exclusiveMax)
}

// LoadGame implements Query.
func (q RawQuery) LoadGame(path string) (ret query.CacheMap, err error) {
	err = errors.New("load game not implemented")
	return
}

// SaveGame implements Query.
func (q RawQuery) SaveGame(path string, dynamicValues query.CacheMap) error {
	return errors.New("load game not implemented")
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

// fix? calls are not consistent about which name they use
// see also BuildKind
func (q *Data) GetKindByName(kind string) (ret *rt.Kind, err error) {
	if k, ok := FindKind(q.Kinds, kind); ok {
		ret = k
	} else if plural, ok := FindPlural(q.Plurals, kind); !ok {
		err = fmt.Errorf("couldnt find kind %q", kind)
	} else if k, ok := FindKind(q.Kinds, plural.Many); ok {
		ret = k
	} else {
		err = fmt.Errorf("couldnt find kind %q or %q", kind, plural.Many)
	}
	return
}

// shortname to id
func (q RawQuery) findFullName(shortname string) (ret string, err error) {
	if a, ok := FindName(q.Names, shortname); !ok {
		err = fmt.Errorf("couldnt find noun with shortname %q", shortname)
	} else {
		ret = a.Noun
	}
	return
}

// fullname to NounData
func (q RawQuery) findNounData(fullname string) (ret NounData, err error) {
	if a, ok := FindNoun(q.Nouns, fullname); !ok {
		err = fmt.Errorf("couldnt find noun with fullname %q", fullname)
	} else {
		ret = a
	}
	return
}

func (q RawQuery) findPattern(name string) (ret PatternData, err error) {
	if a, ok := FindPattern(q.Patterns, name); !ok {
		err = fmt.Errorf("couldnt find pattern %q", name)
	} else {
		ret = a
	}
	return
}

func (q RawQuery) findRelation(name string) (ret *RelativeData, err error) {
	if a, ok := FindRelation(q.Relatives, name); !ok {
		err = fmt.Errorf("couldnt find relation %q", name)
	} else {
		ret = a
	}
	return
}
