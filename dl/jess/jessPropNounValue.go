package jess

import (
	"fmt"
)

func (op *PropertyNounValue) GetBuilder() Builder {
	return &op.builder
}

func TryPropertyNounValue(q JessContext, in InputState,
	accept func(PromisedMatcher), reject func(error),
) {
	var are Are
	if lhs, rhs, ok := are.Split(in); !ok {
		reject(FailedMatch{"property noun phrases require 'is' or 'are'", in})
	} else {
		// the target can be implied
		if propTargetOf, ok := keywordSplit(lhs, keywords.Of); !ok {
			// an implied pronoun; requests the pronoun because we need the kind
			TryImpliedPronoun(q, func(pn Pronoun) {
				pb := PropertyBuilder{Context: q, noun: &pn}
				tryPropertyNounValue(&pb, are, lhs, rhs,
					func(prop Property, val PropertyValue) {
						accept(&PropertyNounValue{
							Property:      prop,
							PropertyNoun:  &pn,
							Are:           are,
							PropertyValue: val,
							builder:       pb,
						})
					}, reject)
			}, reject)
		} else {
			// match a name to a noun ( or generate one )
			// ( interesting to note that inform doesn't allow "kind called" here. )
			TryPropertyNoun(q, propTargetOf.rhs, func(pn PropertyNoun, in InputState) {
				pb := PropertyBuilder{Context: q, noun: pn}
				tryPropertyNounValue(&pb, are, propTargetOf.lhs, rhs,
					func(prop Property, pv PropertyValue) {
						accept(&PropertyNounValue{
							Property:      prop,
							Of:            Words{Matched: propTargetOf.matched},
							PropertyNoun:  pn,
							Are:           are,
							PropertyValue: pv,
							builder:       pb,
						})
					}, reject)
			}, reject)
		}
	}
}

// a super specific function: part of TryPropertyNounValue
func tryPropertyNounValue(pb *PropertyBuilder, isAre Are,
	inProp, inValue InputState,
	accept func(Property, PropertyValue),
	reject func(error),
) {
	// match a property name
	TryPropertyName(pb.Context, inProp, pb.GetKind(), func(prop Property, rest InputState) {
		if rest.Len() != 0 {
			e := fmt.Errorf("trying to define a property for a kind of %q, but didn't recognize %s", pb.GetKind(), rest.DebugString())
			reject(e)
		} else {
			flags := AllowSingular
			if isAre.IsPlural() {
				flags = AllowPlural
			}
			TryPropertyValue(pb.Context, inValue, flags, func(val PropertyValue, fin InputState) {
				if fin.Len() != 0 {
					reject(FailedMatch{"unexpected words after a property definition", rest})
				} else {
					pb.addProperty(prop.fieldName, val.Assignment())
					accept(prop, val)
				}
			}, reject)
		}
	}, reject)
}
