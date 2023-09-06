package query

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

func (q QueryNone) IsDomainActive(name string) (_ bool, _ error) {
	return
}

func (q QueryNone) ActivateDomains(name string) (_, _ []string, err error) {
	err = NotImplemented(q)
	return
}

func (q QueryNone) ReadChecks(actuallyJustThisOne string) (_ []CheckData, _ error) {
	return
}

func (q QueryNone) FieldsOf(kind string) (_ []FieldData, _ error) {
	return
}

func (q QueryNone) KindOfAncestors(kind string) (_ []string, _ error) {
	return
}

func (q QueryNone) KindValues(id string) (_ []ValueData, _ error) {
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

func (q QueryNone) NounValues(id, field string) (_ []ValueData, _ error) {
	return
}

func (q QueryNone) NounsByKind(kind string) (_ []string, _ error) {
	return
}

func (q QueryNone) OppositeOf(kind string) (_ string, _ error) {
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

func (q QueryNone) RulesFor(pat string) (_ []RuleData, _ error) {
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
