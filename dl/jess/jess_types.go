// Code generated by Tapestry; edit at your own risk.
package jess

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// matched, a type of slot.
var Zt_Matched = typeinfo.Slot{
	Name: "matched",
	Markup: map[string]any{
		"comment": []interface{}{"a snippet of matching text", "defined via an interface to allow instances", "to track additional information (ex. db row)"},
	},
}

// holds a single slot.
type Matched_Slot struct{ Value Matched }

// implements typeinfo.Instance for a single slot.
func (*Matched_Slot) TypeInfo() typeinfo.T {
	return &Zt_Matched
}

// holds a slice of slots.
type Matched_Slots []Matched

// implements typeinfo.Instance for a series of slots.
func (*Matched_Slots) TypeInfo() typeinfo.T {
	return &Zt_Matched
}

// implements typeinfo.Repeats
func (op *Matched_Slots) Repeats() bool {
	return len(*op) > 0
}

type Article struct {
	Matched Matched
	Markup  map[string]any
}

// article, a type of flow.
var Zt_Article typeinfo.Flow

// implements typeinfo.Instance
func (*Article) TypeInfo() typeinfo.T {
	return &Zt_Article
}

// implements typeinfo.Markup
func (op *Article) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type article
type Article_Slice []Article

// implements typeinfo.Instance
func (*Article_Slice) TypeInfo() typeinfo.T {
	return &Zt_Article
}

// implements typeinfo.Repeats
func (op *Article_Slice) Repeats() bool {
	return len(*op) > 0
}

// conjunction junction
type CommaAnd struct {
	Matched Matched
	Markup  map[string]any
}

// comma_and, a type of flow.
var Zt_CommaAnd typeinfo.Flow

// implements typeinfo.Instance
func (*CommaAnd) TypeInfo() typeinfo.T {
	return &Zt_CommaAnd
}

// implements typeinfo.Markup
func (op *CommaAnd) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type comma_and
type CommaAnd_Slice []CommaAnd

// implements typeinfo.Instance
func (*CommaAnd_Slice) TypeInfo() typeinfo.T {
	return &Zt_CommaAnd
}

// implements typeinfo.Repeats
func (op *CommaAnd_Slice) Repeats() bool {
	return len(*op) > 0
}

// matches "is" or "are".
type Are struct {
	Matched Matched
	Markup  map[string]any
}

// are, a type of flow.
var Zt_Are typeinfo.Flow

// implements typeinfo.Instance
func (*Are) TypeInfo() typeinfo.T {
	return &Zt_Are
}

// implements typeinfo.Markup
func (op *Are) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type are
type Are_Slice []Are

// implements typeinfo.Instance
func (*Are_Slice) TypeInfo() typeinfo.T {
	return &Zt_Are
}

// implements typeinfo.Repeats
func (op *Are_Slice) Repeats() bool {
	return len(*op) > 0
}

// matches the name of an (existing) trait.
type TraitName struct {
	Matched Matched
	Markup  map[string]any
}

// trait_name, a type of flow.
var Zt_TraitName typeinfo.Flow

// implements typeinfo.Instance
func (*TraitName) TypeInfo() typeinfo.T {
	return &Zt_TraitName
}

// implements typeinfo.Markup
func (op *TraitName) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type trait_name
type TraitName_Slice []TraitName

// implements typeinfo.Instance
func (*TraitName_Slice) TypeInfo() typeinfo.T {
	return &Zt_TraitName
}

// implements typeinfo.Repeats
func (op *TraitName_Slice) Repeats() bool {
	return len(*op) > 0
}

// matches the name of an (existing) kind.
type KindName struct {
	Matched Matched
	Markup  map[string]any
}

// kind_name, a type of flow.
var Zt_KindName typeinfo.Flow

// implements typeinfo.Instance
func (*KindName) TypeInfo() typeinfo.T {
	return &Zt_KindName
}

// implements typeinfo.Markup
func (op *KindName) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type kind_name
type KindName_Slice []KindName

// implements typeinfo.Instance
func (*KindName_Slice) TypeInfo() typeinfo.T {
	return &Zt_KindName
}

// implements typeinfo.Repeats
func (op *KindName_Slice) Repeats() bool {
	return len(*op) > 0
}

// matches at least one trait.
type Traits struct {
	Article          *Article
	TraitName        TraitName
	AdditionalTraits *AdditionalTraits
	Markup           map[string]any
}

// traits, a type of flow.
var Zt_Traits typeinfo.Flow

// implements typeinfo.Instance
func (*Traits) TypeInfo() typeinfo.T {
	return &Zt_Traits
}

// implements typeinfo.Markup
func (op *Traits) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type traits
type Traits_Slice []Traits

// implements typeinfo.Instance
func (*Traits_Slice) TypeInfo() typeinfo.T {
	return &Zt_Traits
}

// implements typeinfo.Repeats
func (op *Traits_Slice) Repeats() bool {
	return len(*op) > 0
}

// matches a trait following another trait
type AdditionalTraits struct {
	CommaAnd *CommaAnd
	Traits   Traits
	Markup   map[string]any
}

// additional_traits, a type of flow.
var Zt_AdditionalTraits typeinfo.Flow

// implements typeinfo.Instance
func (*AdditionalTraits) TypeInfo() typeinfo.T {
	return &Zt_AdditionalTraits
}

// implements typeinfo.Markup
func (op *AdditionalTraits) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type additional_traits
type AdditionalTraits_Slice []AdditionalTraits

// implements typeinfo.Instance
func (*AdditionalTraits_Slice) TypeInfo() typeinfo.T {
	return &Zt_AdditionalTraits
}

// implements typeinfo.Repeats
func (op *AdditionalTraits_Slice) Repeats() bool {
	return len(*op) > 0
}

type Keywords struct {
	Matched Matched
	Markup  map[string]any
}

// keywords, a type of flow.
var Zt_Keywords typeinfo.Flow

// implements typeinfo.Instance
func (*Keywords) TypeInfo() typeinfo.T {
	return &Zt_Keywords
}

// implements typeinfo.Markup
func (op *Keywords) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type keywords
type Keywords_Slice []Keywords

// implements typeinfo.Instance
func (*Keywords_Slice) TypeInfo() typeinfo.T {
	return &Zt_Keywords
}

// implements typeinfo.Repeats
func (op *Keywords_Slice) Repeats() bool {
	return len(*op) > 0
}

type MacroName struct {
	Matched Matched
	Macro   Macro
	Markup  map[string]any
}

// macro_name, a type of flow.
var Zt_MacroName typeinfo.Flow

// implements typeinfo.Instance
func (*MacroName) TypeInfo() typeinfo.T {
	return &Zt_MacroName
}

// implements typeinfo.Markup
func (op *MacroName) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type macro_name
type MacroName_Slice []MacroName

// implements typeinfo.Instance
func (*MacroName_Slice) TypeInfo() typeinfo.T {
	return &Zt_MacroName
}

// implements typeinfo.Repeats
func (op *MacroName_Slice) Repeats() bool {
	return len(*op) > 0
}

type KindsAreTraits struct {
	Article  *Article
	KindName KindName
	Are      Are
	Usually  MacroName
	Traits   Traits
	Markup   map[string]any
}

// kinds_are_traits, a type of flow.
var Zt_KindsAreTraits typeinfo.Flow

// implements typeinfo.Instance
func (*KindsAreTraits) TypeInfo() typeinfo.T {
	return &Zt_KindsAreTraits
}

// implements typeinfo.Markup
func (op *KindsAreTraits) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// holds a slice of type kinds_are_traits
type KindsAreTraits_Slice []KindsAreTraits

// implements typeinfo.Instance
func (*KindsAreTraits_Slice) TypeInfo() typeinfo.T {
	return &Zt_KindsAreTraits
}

// implements typeinfo.Repeats
func (op *KindsAreTraits_Slice) Repeats() bool {
	return len(*op) > 0
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name:       "jess",
	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Matched,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_Article,
	&Zt_CommaAnd,
	&Zt_Are,
	&Zt_TraitName,
	&Zt_KindName,
	&Zt_Traits,
	&Zt_AdditionalTraits,
	&Zt_Keywords,
	&Zt_MacroName,
	&Zt_KindsAreTraits,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	508023169458945308:   (*AdditionalTraits)(nil), /* AdditionalTraits commaAnd:traits: */
	1887918947148326916:  (*AdditionalTraits)(nil), /* AdditionalTraits traits: */
	503552734697422485:   (*Are)(nil),              /* Are: */
	8854300316672007225:  (*Article)(nil),          /* Article: */
	4230553755039810705:  (*CommaAnd)(nil),         /* CommaAnd: */
	3748071630827580029:  (*Keywords)(nil),         /* Keywords: */
	11329968067422663542: (*KindName)(nil),         /* KindName: */
	10676540530754351752: (*KindsAreTraits)(nil),   /* KindsAreTraits article:kindName:are:usually:traits: */
	13876633153604694564: (*KindsAreTraits)(nil),   /* KindsAreTraits kindName:are:usually:traits: */
	15972029076576488422: (*MacroName)(nil),        /* MacroName: */
	13073468751382026622: (*TraitName)(nil),        /* TraitName: */
	12266123994356951879: (*Traits)(nil),           /* Traits article:traitName: */
	9462642200558625635:  (*Traits)(nil),           /* Traits article:traitName:additionalTraits: */
	10912000954309087975: (*Traits)(nil),           /* Traits traitName: */
	11584307469651821251: (*Traits)(nil),           /* Traits traitName:additionalTraits: */
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_Article = typeinfo.Flow{
		Name: "article",
		Lede: "article",
		Terms: []typeinfo.Term{{
			Name: "matched",
			Type: &Zt_Matched,
		}},
		Markup: map[string]any{
			"comment": []interface{}{"one of a predefined set of determiners:", "the, a, some, etc.", "see 'count' for leading numbers"},
		},
	}
	Zt_CommaAnd = typeinfo.Flow{
		Name: "comma_and",
		Lede: "comma_and",
		Terms: []typeinfo.Term{{
			Name: "matched",
			Type: &Zt_Matched,
		}},
		Markup: map[string]any{
			"comment": "conjunction junction",
		},
	}
	Zt_Are = typeinfo.Flow{
		Name: "are",
		Lede: "are",
		Terms: []typeinfo.Term{{
			Name: "matched",
			Type: &Zt_Matched,
		}},
		Markup: map[string]any{
			"comment": "matches \"is\" or \"are\".",
		},
	}
	Zt_TraitName = typeinfo.Flow{
		Name: "trait_name",
		Lede: "trait_name",
		Terms: []typeinfo.Term{{
			Name: "matched",
			Type: &Zt_Matched,
		}},
		Markup: map[string]any{
			"comment": "matches the name of an (existing) trait.",
		},
	}
	Zt_KindName = typeinfo.Flow{
		Name: "kind_name",
		Lede: "kind_name",
		Terms: []typeinfo.Term{{
			Name: "matched",
			Markup: map[string]any{
				"comment": "the matched type is always a span.",
			},
			Type: &Zt_Matched,
		}},
		Markup: map[string]any{
			"comment": "matches the name of an (existing) kind.",
		},
	}
	Zt_Traits = typeinfo.Flow{
		Name: "traits",
		Lede: "traits",
		Terms: []typeinfo.Term{{
			Name:     "article",
			Label:    "article",
			Optional: true,
			Markup: map[string]any{
				"comment": []interface{}{"while an article can precede every trait", "it doesn't influence which trait gets matched."},
			},
			Type: &Zt_Article,
		}, {
			Name:  "trait_name",
			Label: "trait_name",
			Type:  &Zt_TraitName,
		}, {
			Name:     "additional_traits",
			Label:    "additional_traits",
			Optional: true,
			Type:     &Zt_AdditionalTraits,
		}},
		Markup: map[string]any{
			"comment": "matches at least one trait.",
		},
	}
	Zt_AdditionalTraits = typeinfo.Flow{
		Name: "additional_traits",
		Lede: "additional_traits",
		Terms: []typeinfo.Term{{
			Name:     "comma_and",
			Label:    "comma_and",
			Optional: true,
			Type:     &Zt_CommaAnd,
		}, {
			Name:  "traits",
			Label: "traits",
			Type:  &Zt_Traits,
		}},
		Markup: map[string]any{
			"comment": "matches a trait following another trait",
		},
	}
	Zt_Keywords = typeinfo.Flow{
		Name: "keywords",
		Lede: "keywords",
		Terms: []typeinfo.Term{{
			Name: "matched",
			Type: &Zt_Matched,
		}},
		Markup: map[string]any{
			"comment": []interface{}{"matches one or more predefined words", "the specific words are specified via metadata", "on the term where this flow is declared."},
		},
	}
	Zt_MacroName = typeinfo.Flow{
		Name: "macro_name",
		Lede: "macro_name",
		Terms: []typeinfo.Term{{
			Name: "matched",
			Type: &Zt_Matched,
		}, {
			Name:    "macro",
			Label:   "macro",
			Private: true,
		}},
		Markup: map[string]any{
			"comment": []interface{}{"matches one or more predefined words", "and returns a macro. like keywords", "the phrase and macro are on defined via metadata"},
		},
	}
	Zt_KindsAreTraits = typeinfo.Flow{
		Name: "kinds_are_traits",
		Lede: "kinds_are_traits",
		Terms: []typeinfo.Term{{
			Name:     "article",
			Label:    "article",
			Optional: true,
			Type:     &Zt_Article,
		}, {
			Name:  "kind_name",
			Label: "kind_name",
			Type:  &Zt_KindName,
		}, {
			Name:  "are",
			Label: "are",
			Type:  &Zt_Are,
		}, {
			Name:  "usually",
			Label: "usually",
			Markup: map[string]any{
				"macro":  "implies",
				"phrase": "usually",
			},
			Type: &Zt_MacroName,
		}, {
			Name:  "traits",
			Label: "traits",
			Type:  &Zt_Traits,
		}},
		Markup: map[string]any{
			"comment": []interface{}{"assigns default traits to a kind.", "[the] <kind> are \"usually\" <traits>"},
		},
	}
}
