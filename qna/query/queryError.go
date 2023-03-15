package query

// NotImplemented - implements Query by returning a NotImplementedError error for every method.
type NotImplemented string

// NotImplementedError - generic error used returned by NotImplemented
type NotImplementedError string

func (e NotImplementedError) Error() string {
	return string(e)
}

func (q NotImplemented) IsDomainActive(name string) (ret bool, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) ActivateDomain(name string) (ret string, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) ReadChecks(actuallyJustThisOne string) (ret []CheckData, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) FieldsOf(kind string) (ret []FieldData, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) NounInfo(name string) (ret NounInfo, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) NounIsNamed(id, name string) (ret bool, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) NounName(id string) (ret string, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) NounValue(id, field string) (ret []byte, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) NounsByKind(kind string) (ret []string, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) PluralToSingular(plural string) (ret string, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) PluralFromSingular(singular string) (ret string, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) PatternLabels(pat string) (ret []string, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) RulesFor(pat, target string) (ret []Rules, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) ReciprocalsOf(rel, id string) (ret []string, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) RelativesOf(rel, id string) (ret []string, err error) {
	err = NotImplementedError(q)
	return
}

func (q NotImplemented) Relate(rel, noun, otherNoun string) error {
	return NotImplementedError(q)
}
