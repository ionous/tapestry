package jess

import (
	"git.sr.ht/~ionous/tapestry/support/grok"
)

func (op *Are) Match(q Query, input *InputState) (okay bool) {
	if width := input.MatchWord(grok.Keyword.Are, grok.Keyword.Is); width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *CommaAnd) Match(q Query, input *InputState) (okay bool) {
	if sep, e := grok.CommaAnd(input.Words()); e != nil {
		q.error("comma and", e)
	} else {
		width := sep.Len()
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *Keywords) Match(q Query, input *InputState, meta map[string]any) (okay bool) {
	if words, ok := meta["keywords"]; !ok {
		q.log("missing keyword metadata")
	} else if width := match(input, words); width < 0 {
		q.log("invalid keyword metadata")
	} else if width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *MacroName) Match(q Query, input *InputState, meta map[string]any) (okay bool) {
	if words, ok := meta["phrase"].(string); !ok || len(words) == 0 {
		q.log("missing keyword metadata")
	} else if phrase, e := grok.MakeSpan(words); e != nil {
		q.error("macro", e)
	} else if m, e := q.g.FindMacro(phrase); e != nil {
		q.error("find macro", e)
	} else if grok.HasPrefix(input.Words(), phrase) {
		width := len(phrase)
		op.Macro = m
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func match(input *InputState, meta any) (ret int) {
	switch ws := meta.(type) {
	case string:
		ret = input.MatchWord(grok.Hash(ws))
	case []string:
		ret = input.MatchWord(grok.Hashes(ws)...)
	default:
		ret = -1
	}
	return
}
