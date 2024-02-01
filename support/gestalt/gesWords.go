package gestalt

import "git.sr.ht/~ionous/tapestry/support/grok"

type Are struct{}

func (*Are) Match(q Query, cs []InputState) (ret []InputState) {
	return matchAll(cs, grok.Keyword.Is, grok.Keyword.Are)
}

// matches a specified string and advances the input
type Words struct {
	Str string
	matcher
}

func (op *Words) Match(q Query, cs []InputState) (ret []InputState) {
	if match, e := op.makeSpan(op.Str); e == nil {
		for _, in := range cs {
			if grok.HasPrefix(in.Words(), match) {
				rest := in.Next(len(match))
				ret = append(ret, rest)
			}
		}
	}
	return
}

type matcher struct {
	cache grok.Span
	err   error
}

// caches passed string, doesnt check to see if the string changes.
func (m *matcher) makeSpan(str string) (grok.Span, error) {
	if m.cache == nil && m.err == nil && len(str) > 0 {
		m.cache, m.err = grok.MakeSpan(str)
	}
	return m.cache, m.err
}

func matchAll(cs []InputState, choices ...uint64) (ret []InputState) {
	for _, in := range cs {
		if matchWord(in.Words(), choices...) > 0 {
			ret = append(ret, in.Next(1))
		}
	}
	return
}

func matchWord(ws []grok.Word, choices ...uint64) (width int) {
	if len(ws) > 0 {
		w := ws[0]
		for _, opt := range choices {
			if w.Equals(opt) {
				width = 1
				break
			}
		}
	}
	return
}
