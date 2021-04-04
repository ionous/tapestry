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

// changes the bounds of its first scanner in response to the results of its last scanner.
func retarget(s ...parser.Scanner) parser.Scanner {
	return &parser.Target{s}
}

func noun(kinds ...string) parser.Scanner {
	var fs parser.Filters
	for _, k := range kinds {
		fs = append(fs, &parser.HasClass{k})
	}
	return &parser.Noun{fs}
}
func nouns(kinds ...string) parser.Scanner {
	var fs parser.Filters
	for _, k := range kinds {
		fs = append(fs, &parser.HasClass{k})
	}
	return &parser.Multi{fs}
}

// note: we use things to exclude directions
func thing() parser.Scanner {
	return noun("things")
}

func things() parser.Scanner {
	return nouns("things")
}

func words(s string) parser.Scanner {
	return parser.Words(s)
}

func act(s string) *parser.Action {
	return &parser.Action{s}
}
