package jess

import "fmt"

// see: TryPropertyNounValue
func (op *PropertyNounValue) PromiseMatcher() PromiseMatcher {
	return op
}

func TryPropertyNounValue(q JessContext, in InputState,
	accept func(PromiseMatcher), reject func(error),
) {
	var are Are
	if descValueIs, ok := are.Split(in); !ok {
		reject(FailedMatch{"property noun phrases require 'is' or 'are'", in})
	} else {
		// the target can be implied
		if propTargetOf, ok := keywordSplit(descValueIs.lhs, keywords.Of); !ok {
			RequestPronoun(q, func(an ActualNoun) {
				GeneratePropertyNounValue(q,
					an, are,
					descValueIs.lhs,
					descValueIs.rhs,
					func(prop Property, val PropertyValue) {
						accept(&PropertyNounValue{
							Property:      prop,
							Are:           are,
							PropertyValue: val,
						})
					}, reject)
			}, reject)
		} else {
			// match a name to a noun ( or generate one )
			// ( interesting to note that inform doesn't allow "kind called" here. )
			TryNamedNoun(q, propTargetOf.rhs, func(nn NamedNoun, an ActualNoun, in InputState) {
				GeneratePropertyNounValue(q,
					an, are,
					propTargetOf.lhs,
					descValueIs.rhs,
					func(prop Property, val PropertyValue) {
						accept(&PropertyNounValue{
							Property:      prop,
							Of:            Words{Matched: propTargetOf.matched},
							NamedNoun:     nn,
							Are:           are,
							PropertyValue: val,
						})
					}, reject)
			}, reject)
		}
	}
}

// a super specific function: part of TryPropertyNounValue
func GeneratePropertyNounValue(q JessContext,
	an ActualNoun, isAre Are,
	inProp, inValue InputState,
	accept func(Property, PropertyValue),
	reject func(error),
) {
	// match a property name
	TryProperty(q, inProp, an.Kind, func(prop Property, rest InputState) {
		if rest.Len() != 0 {
			e := fmt.Errorf("trying to define a property for a noun %q of kind %q, but didn't recognize %s", an.Name, an.Kind, rest.DebugString())
			reject(e)
		} else {
			generatePropertyValue(q, inValue,
				an.Name, prop.fieldName, isAre.IsPlural(),
				func(val PropertyValue, rest InputState) {
					if rest.Len() != 0 {
						reject(FailedMatch{"unexpected words trailing a property definition", rest})
					} else {
						accept(prop, val)
					}
				}, reject)
		}
	}, reject)
}
