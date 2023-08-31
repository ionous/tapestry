package play

import "git.sr.ht/~ionous/tapestry/parser"

// adapt from Playtime to parser context
type parserContext Playtime

func (pt *parserContext) GetBounds(who, where string) (parser.Bounds, error) {
	p := (*Playtime)(pt)
	return p.survey.GetBounds(who, where)
}

func (pt *parserContext) IsPlural(word string) bool {
	p := (*Playtime)(pt)
	pl := p.Runtime.SingularOf(word)
	return len(pl) > 0 && pl != word
}
