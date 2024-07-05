package query

import "git.sr.ht/~ionous/tapestry/rt"

// QueryNone - implements Query by returning empty results for all reads,
// and the NotImplemented error for mutating methods.
type QueryNone string

// verify that query none implements every method
var _ = Query(QueryNone(""))

// NotImplemented - generic error used returned by QueryNone
type NotImplemented string

func (e NotImplemented) Error() string {
	return string(e)
}

func (q QueryNone) Close() {
}

func (q QueryNone) GetKindByName(rawName string) (_ *rt.Kind, err error) {
	err = NotImplemented(q)
	return
}

func (q QueryNone) IsDomainActive(name string) (_ bool, _ error) {
	return
}

func (q QueryNone) ActivateDomains(name string) (_, _ []string, err error) {
	err = NotImplemented(q)
	return
}

func (q QueryNone) KindOfAncestors(kind string) (_ []string, _ error) {
	return
}

func (q QueryNone) NounInfo(name string) (_ NounInfo, _ error) {
	return
}

func (q QueryNone) NounName(id string) (_ string, _ error) {
	return
}

func (q QueryNone) NounNames(id string) (_ []string, _ error) {
	return
}

func (q QueryNone) NounValue(id, field string) (_ rt.Assignment, _ error) {
	return
}

func (q QueryNone) NounsByKind(kind string) (_ []string, _ error) {
	return
}

func (q QueryNone) PluralToSingular(plural string) (_ string, _ error) {
	return
}

func (q QueryNone) PluralFromSingular(singular string) (_ string, _ error) {
	return
}

func (q QueryNone) PatternLabels(pat string) (_ []string, _ error) {
	return
}

func (q QueryNone) RulesFor(pat string) (_ RuleSet, _ error) {
	return
}

func (q QueryNone) ReciprocalsOf(rel, id string) (_ []string, _ error) {
	return
}

func (q QueryNone) RelativesOf(rel, id string) (_ []string, _ error) {
	return
}

func (q QueryNone) Relate(rel, noun, otherNoun string) error {
	return NotImplemented(q)
}

// Random implements Query.
func (q QueryNone) Random(inclusiveMin int, exclusiveMax int) int {
	return inclusiveMin
}

// LoadGame implements Query.
func (q QueryNone) LoadGame(path string) (ret CacheMap, err error) {
	err = NotImplemented(q)
	return
}

// SaveGame implements Query.
func (q QueryNone) SaveGame(path string, dynamicValues CacheMap) error {
	return NotImplemented(q)
}
