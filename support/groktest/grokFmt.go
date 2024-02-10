package groktest

import "git.sr.ht/~ionous/tapestry/support/grok"

func ResultMap(in grok.Results) map[string]any {
	m := make(map[string]any)
	nounsIntoMap(m, "primary", in.Primary)
	nounsIntoMap(m, "secondary", in.Secondary)
	macroIntoMap(m, "macro", in.Macro)
	return m
}

func traitSetMap(ts grok.TraitSet) map[string]any {
	m := make(map[string]any)
	matchesIntoMap(m, "traits", ts.Traits)
	matchIntoMap(m, "kind", ts.Kind)
	return m
}

func nounsIntoMap(m map[string]any, field string, ns []grok.Name) {
	if len(ns) > 0 {
		out := make([]map[string]any, len(ns))
		for i, n := range ns {
			out[i] = nounToMap(n)
		}
		m[field] = out
	}
}

func nounToMap(n grok.Name) map[string]any {
	m := make(map[string]any)
	matchIntoMap(m, "name", n.Span)
	matchIntoMap(m, "det", n.Article.Matched)
	matchesIntoMap(m, "traits", n.Traits)
	matchesIntoMap(m, "kinds", n.Kinds)
	if n.Exact {
		m["exact"] = true
	}
	if cnt := n.Article.Count; cnt > 0 {
		m["count"] = cnt
	}
	return m
}

func matchesIntoMap(m map[string]any, field string, ws []grok.Matched) {
	if cnt := len(ws); cnt > 0 {
		out := make([]string, cnt)
		for i, w := range ws {
			out[i] = w.String()
		}
		m[field] = out
	}
	return
}

func matchIntoMap(m map[string]any, field string, match grok.Matched) {
	if match != nil && match.NumWords() > 0 {
		m[field] = match.String()
	}
}

func macroIntoMap(m map[string]any, field string, macro grok.Macro) {
	if macro.Matched != nil && macro.Matched.NumWords() > 0 {
		m[field] = macro.Name
	}
}
