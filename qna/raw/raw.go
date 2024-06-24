// a non-sql data format for use by the qna runtime
package raw

import "git.sr.ht/~ionous/tapestry/qna/query"

type Data struct {
	Scenes    []string   // in order with root ( tapestry ) first
	Plurals   []Plural   // pairs of one, many; sorted by one
	Kinds     []KindData // sorted by .Kind
	Names     []NounName // sorted by .Name
	Nouns     []NounData
	Patterns  []PatternData
	Relatives []RelativeData // sorted by .Relation
}

type KindData struct {
	Id        int
	Domain    string
	Kind      string
	Ancestors []string // includes itself
	Fields    []query.FieldData
}

type NounName struct {
	Id   int
	Name string
}

type NounData struct {
	Id     int // mdl_noun
	Domain string
	Noun   string // unique id
	Kind   string // or would id be better?
	// Name list? or best name maybe?
	Values []query.ValueData
}

type PatternData struct {
	Id      int // mdl_pat
	Pattern string
	Labels  []string
	Result  string
	Rules   []query.RuleData
}

type RelativeData struct {
	Id       int    // mdl_rel
	Relation string // a kind
	Pairs    []Pair // can be empty
}

type Plural struct {
	One, Other string
}

type Pair struct {
	One, Other int // mdl_noun
}
