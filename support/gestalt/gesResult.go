package gestalt

import "git.sr.ht/~ionous/tapestry/support/grok"

// the results carry match ( instead of just a range or something )
// so that the matcher can store custom data ( ex. db row )
type Result struct {
	Result ResultType
	Match  grok.Match
}

type ResultType int

//go:generate stringer -type=ResultType
const (
	None ResultType = iota
	TypeArticle
	TypeExactName
	TypeKind
	TypeTrait
)

func Reduce(in InputState) (ret grok.Results) {
	var name grok.Name
	for _, res := range in.res {
		m := res.Match
		switch res.Result {
		case TypeArticle:
			name.Article = m.(grok.Article)

		case TypeExactName:
			name.Span = m.(grok.Span)
			name.Exact = true

		case TypeKind:
			name.Kinds = append(name.Kinds, m)

		case TypeTrait:
			name.Traits = append(name.Traits, m)
		}
	}

	ret.Primary = append(ret.Primary, name)
	return
}
