// a non-sql data format for use by the qna runtime
package raw

import (
	"git.sr.ht/~ionous/tapestry/dl/grammar"
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
	Scenes    []string            // in order with root ( tapestry ) first
	Plurals   []Plural            // pairs of one, many; sorted by one
	Kinds     []rt.Kind           // sorted by .Kind
	Names     []NounName          // sorted by .Name
	Nouns     []NounData          // sorted by .Noun
	Patterns  []PatternData       // sorted by .Pattern
	Relatives []RelativeData      // sorted by .Relation
	Grammar   []grammar.Directive // sorted by .Name
}

type SceneData struct {
	Scene    string
	Requires []string `json:",omitempty"`
}

type NounName struct {
	Name, Noun string // shortname, fullname
}

type NounData struct {
	Id         int // mdl_noun
	Domain     string
	Noun       string       // full, unique name
	Kind       string       // or would id be better?
	CommonName string       // author defined name
	Aliases    []string     // alpha order for parser
	Values     []EvalData   `json:",omitempty"` // sparse field values, sorted by field name
	Records    []RecordData `json:",omitempty"`
}

type EvalData struct {
	Field string
	Value rt.Assignment
}

// record data is json serialized ( via pack/unpack )
// because otherwise gobbing would need a way to store variant values
// package json says []byte encodes as a base64-encoded string,
type RecordData struct {
	Field  string
	Packed []byte // json serialized for now
}

type PatternData struct {
	Id        int // mdl_pat
	Pattern   string
	Labels    []string  // the last value is the result, even if blank
	Rules     []rt.Rule `json:",omitempty"`
	UpdateAll bool
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
	One, Other string // noun fullnames
}
