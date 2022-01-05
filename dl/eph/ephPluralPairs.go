package eph

// one plural word can map to multiple single words
type PluralPairs struct {
	// we dont expect there will be huge numbers of plurals,
	// so unsorted arrays should be fine.
	plural, singular []string
}

func (ps PluralPairs) FindPlural(singular string) (ret string, okay bool) {
	if i, ok := find(singular, ps.singular); ok {
		ret, okay = ps.plural[i], true
	}
	return
}
func (ps PluralPairs) FindSingular(plural string) (ret string, okay bool) {
	if i, ok := find(plural, ps.plural); ok {
		ret, okay = ps.singular[i], true
	}
	return
}
func (ps *PluralPairs) AddPair(plural, singular string) (okay bool) {
	if len(plural) > 0 && len(singular) > 0 {
		// is the pairing unique?
		_, havep := find(plural, ps.plural)
		i, haves := find(singular, ps.singular)
		if !havep && (!haves || plural != ps.plural[i]) {
			ps.plural = append(ps.plural, plural)
			ps.singular = append(ps.singular, singular)
			okay = true
		}
	}
	return
}
func find(str string, strs []string) (ret int, okay bool) {
	for i, s := range strs {
		if s == str {
			ret, okay = i, true
			break
		}
	}
	return
}
