// a non-sql data format for use by the qna runtime
package raw

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
)

// cardinality
const (
	ONE_TO_ONE   = "one_one"
	ONE_TO_MANY  = "one_any"
	MANY_TO_ONE  = "any_one"
	MANY_TO_MANY = "any_any"
)

type Data struct {
	Scenes    []string       // in order with root ( tapestry ) first
	Plurals   []Plural       // pairs of one, many; sorted by one
	Kinds     []KindData     // sorted by .Kind
	Names     []NounName     // sorted by .Name
	Nouns     []NounData     // sorted by .Noun
	Patterns  []PatternData  // sorted by .Pattern
	Relatives []RelativeData // sorted by .Relation
	Grammar   []Grammar      // sorted by .Name
}

type KindData struct {
	Id        int
	Domain    string
	Kind      string
	Ancestors []string    // includes itself
	Fields    []FieldData `json:",omitempty"`
}

type FieldData struct {
	Name     string
	Affinity affine.Affinity
	Class    string `json:",omitempty"`
	Init     Bytes  `json:",omitempty"`
}

type NounName struct {
	Name, Noun string // shortname, fullname
}

type NounData struct {
	Id         int // mdl_noun
	Domain     string
	Noun       string      // full, unique name
	Kind       string      // or would id be better?
	CommonName string      `json:",omitempty"` // author defined name
	Aliases    []string    // alpha order for parser
	Values     []ValueData `json:",omitempty"` // sorted by field
}

type ValueData struct {
	Field string
	Path  string `json:",omitempty"`
	Value Bytes  `json:",omitempty"` // a serialized assignment or literal
}

type RuleData struct {
	Name    string
	Stop    bool
	Jump    rt.Jump
	Updates bool
	Prog    Bytes `json:",omitempty"` // a serialized rt.Execute_Slice
}

type PatternData struct {
	Id      int // mdl_pat
	Pattern string
	Labels  []string   // the last value is the result, even if blank
	Rules   []RuleData `json:",omitempty"`
}

type RelativeData struct {
	Id                 int    // mdl_rel
	Relation           string // a kind
	OneKind, OtherKind string // primary and secondary types
	Cardinality        string // ( these also are recorded in the kind data )
	Pairs              []Pair `json:",omitempty"`
}

type Plural struct {
	One, Many string
}

type Pair struct {
	One, Other string // noun fullname
}

type Grammar struct {
	Name string
	Prog Bytes `json:",omitempty"` // a grammar.Directive
}
