package groktest

import "git.sr.ht/~ionous/tapestry/support/grok"

func resultMap(in grok.Results) map[string]any {
	m := make(map[string]any)
	nounsIntoMap(m, "sources", in.Sources)
	nounsIntoMap(m, "targets", in.Targets)
	macroIntoMap(m, "macro", in.Macro)
	return m
}

func traitSetMap(ts grok.TraitSet) map[string]any {
	m := make(map[string]any)
	matchesIntoMap(m, "traits", ts.Traits)
	matchIntoMap(m, "kind", ts.Kind)
	return m
}

func nounsIntoMap(m map[string]any, field string, ns []grok.Noun) {
	if len(ns) > 0 {
		out := make([]map[string]any, len(ns))
		for i, n := range ns {
			out[i] = nounToMap(n)
		}
		m[field] = out
	}
}

func nounToMap(n grok.Noun) map[string]any {
	m := make(map[string]any)
	matchIntoMap(m, "name", n.Name)
	matchIntoMap(m, "det", n.Det)
	matchesIntoMap(m, "traits", n.Traits)
	matchesIntoMap(m, "kinds", n.Kinds)
	if n.Exact {
		m["exact"] = true
	}
	return m
}

func matchesIntoMap(m map[string]any, field string, ws []grok.Match) {
	if cnt := len(ws); cnt > 0 {
		out := make([]string, cnt)
		for i, w := range ws {
			out[i] = w.String()
		}
		m[field] = out
	}
	return
}

func matchIntoMap(m map[string]any, field string, match grok.Match) {
	if match != nil && match.NumWords() > 0 {
		m[field] = match.String()
	}
}

func macroIntoMap(m map[string]any, field string, macro grok.MacroInfo) {
	if macro.Match != nil && macro.Match.NumWords() > 0 {
		m[field] = macro.Name
	}
}
