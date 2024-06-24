// a non-sql data format for use by the qna runtime
package raw

import "git.sr.ht/~ionous/tapestry/qna/query"

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
	Ancestors []string // includes itself
	Fields    []query.FieldData
}

type NounName struct {
	Name, Noun string // shortname, fullname
}

type NounData struct {
	Id      int // mdl_noun
	Domain  string
	Noun    string            // fullname
	Kind    string            // or would id be better?
	Aliases []string          // the friendly name is first, followed by specification order
	Values  []query.ValueData // sorted by field
}

type PatternData struct {
	Id      int // mdl_pat
	Pattern string
	Labels  []string // the last value is the result, even if blank
	Rules   []query.RuleData
}

type RelativeData struct {
	Id       int    // mdl_rel
	Relation string // a kind
	Pairs    []Pair // can be empty
}

type Plural struct {
	One, Many string
}

type Pair struct {
	One, Other string // noun fullname
}

type Grammar struct {
	Name string
	Prog query.Bytes `json:",omitempty"` // a grammar.Directive
}
