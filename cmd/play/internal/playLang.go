package internal

import "git.sr.ht/~ionous/iffy/parser"

func anyOf(s ...parser.Scanner) (ret parser.Scanner) {
	if len(s) == 1 {
		ret = s[0]
	} else {
		ret = &parser.AnyOf{s}
	}
	return
}

func allOf(s ...parser.Scanner) (ret parser.Scanner) {
	if len(s) == 1 {
		ret = s[0]
	} else {
		ret = &parser.AllOf{s}
	}
	return
}

func noun(f ...parser.Filter) parser.Scanner {
	return &parser.Noun{f}
}
func nouns(f ...parser.Filter) parser.Scanner {
	return &parser.Multi{f}
}

// note: we use things to exclude directions
func thing() parser.Scanner {
	return noun(&parser.HasClass{"things"})
}

func things() parser.Scanner {
	return nouns(&parser.HasClass{"things"})
}

func words(s string) parser.Scanner {
	return parser.Words(s)
}

func act(s string) *parser.Action {
	return &parser.Action{s}
}
