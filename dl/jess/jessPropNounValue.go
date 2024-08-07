package jess

import (
	"fmt"
)

// see: TryPropertyNounValue
func (op *PropertyNounValue) ParallelMatcher() ParallelMatcher {
	return op
}

// `The description of the pen is "mightier than the sword.`
func TryPropertyNounValue(q JessContext, in InputState,
	accept func(ParallelMatcher), reject func(error),
) {
	// the word "of" splits the property from the noun.
	if propTarget, ok := keywordSplit(in, keywords.Of); !ok {
		reject(FailedMatch{"property noun phrases require the word 'of'", in})
	} else if nameValue, ok := keywordSplit(propTarget.rhs, keywords.Is, keywords.Are); !ok {
		reject(FailedMatch{"property noun phrases require 'is' or 'are'", propTarget.rhs})
	} else {
		// match a name to a noun ( or generate one )
		// ( interesting to note that inform doesn't allow "kind called" here. )
		TryNamedNoun(q, nameValue.lhs, func(nn NamedNoun, an ActualNoun, in InputState) {
			// match a property name
			TryProperty(q, propTarget.lhs, an.Kind, func(prop Property, rest InputState) {
				if rest.Len() != 0 {
					e := fmt.Errorf("It seems you were trying to define a property for a noun %q of kind %q, but i didn't recognize %s", an.Name, an.Kind, rest.DebugString())
					reject(e)
				} else {
					are := Are{Matched: nameValue.matched}
					GeneratePropertyValue(q, nameValue.rhs,
						an.Name, prop.fieldName, are.IsPlural(),
						func(val PropertyValue, final InputState) {
							if final.Len() != 0 {
								reject(FailedMatch{"unexpected words trailing a property definition", final})
							} else {
								accept(&PropertyNounValue{
									Property:      prop,
									Of:            Words{Matched: propTarget.matched},
									NamedNoun:     nn,
									Are:           are,
									PropertyValue: val,
								})
							}
						}, reject)
				}
			}, reject)
		}, reject)
	}

}
