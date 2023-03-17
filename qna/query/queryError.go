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

func (q QueryNone) IsDomainActive(name string) (ret bool, err error) {
	return
}

func (q QueryNone) ActivateDomain(name string) (ret string, err error) {
	err = NotImplemented(q)
	return
}

func (q QueryNone) ReadChecks(actuallyJustThisOne string) (ret []CheckData, err error) {
	return
}

func (q QueryNone) FieldsOf(kind string) (ret []FieldData, err error) {
	return
}

func (q QueryNone) KindOfAncestors(kind string) (ret []string, err error) {
	return
}

func (q QueryNone) NounInfo(name string) (ret NounInfo, err error) {
	return
}

func (q QueryNone) NounIsNamed(id, name string) (ret bool, err error) {
	return
}

func (q QueryNone) NounName(id string) (ret string, err error) {
	return
}

func (q QueryNone) NounValue(id, field string) (ret []byte, err error) {
	return
}

func (q QueryNone) NounsByKind(kind string) (ret []string, err error) {
	return
}

func (q QueryNone) PluralToSingular(plural string) (ret string, err error) {
	return
}

func (q QueryNone) PluralFromSingular(singular string) (ret string, err error) {
	return
}

func (q QueryNone) PatternLabels(pat string) (ret []string, err error) {
	return
}

func (q QueryNone) RulesFor(pat, target string) (ret []Rules, err error) {
	return
}

func (q QueryNone) ReciprocalsOf(rel, id string) (ret []string, err error) {
	return
}

func (q QueryNone) RelativesOf(rel, id string) (ret []string, err error) {
	return
}

func (q QueryNone) Relate(rel, noun, otherNoun string) error {
	return NotImplemented(q)
}
