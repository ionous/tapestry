package parser

import (
	"bytes"
	"strconv"

	"git.sr.ht/~ionous/iffy/ident"
)

// ResultList contains multiple results. Its methods help tease out its contents.
// Most often when a parsing succeeds, it will return a ResultList and
// the .Last() element of the list will be an Action
type ResultList struct {
	list  []Result
	count int
}

// addResult to the list, updating the number of words matched.
func (rs *ResultList) addResult(r Result) {
	if rl, ok := r.(*ResultList); ok {
		rs.list = append(rs.list, rl.list...)
		rs.count += rl.count
	} else {
		rs.list = append(rs.list, r)
		rs.count += r.WordsMatched()
	}
}

// WordsMatched returns the number of words matched.
func (rs *ResultList) WordsMatched() int {
	return rs.count
}

func (rs *ResultList) Results() []Result {
	return rs.list
}

// Last result in the list, true if the list was not empty.
// Generally, when the parser succeeds, this is an Action.
func (rs *ResultList) Last() (ret Result, okay bool) {
	if cnt := len(rs.list); cnt > 0 {
		ret, okay = rs.list[cnt-1], true
	}
	return
}

// Objects -- all nouns used by this result.
func (rs *ResultList) Objects() (ret []ident.Id) {
	for _, r := range rs.list {
		switch k := r.(type) {
		case ResolvedNoun:
			n := k.NounInstance
			ret = append(ret, n.Id())
		case ResolvedMulti:
			for _, n := range k.Nouns {
				ret = append(ret, n.Id())
			}
		}
	}
	return
}

func (rs *ResultList) PrettyObjects() string {
	return PrettyIds(rs.Objects())
}

func (rs *ResultList) String() string {
	var b bytes.Buffer
	b.WriteString("Results(")
	b.WriteString(strconv.Itoa(len(rs.list)))
	b.WriteString("): ")
	for i, res := range rs.list {
		if i > 0 {
			b.WriteString(", ")
		}
		if str, ok := res.(interface{ String() string }); !ok {
			b.WriteString("???")
		} else {
			b.WriteString(str.String())
		}
	}
	return b.String()
}

// PrettyIds - convert a ids to a single comma separated string
func PrettyIds(ids []ident.Id) (ret string) {
	var buf bytes.Buffer
	for i, id := range ids {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(id.String())
	}
	return buf.String()
}
