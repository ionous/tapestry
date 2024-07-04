package compact

import (
	"unicode"
	"unicode/utf8"
)

// Common markup keys in tapestry compact (if) files.
// the comment marker was chosen to be distinctive and to sort in front of all other keys
// (ie. when using go's json encoder. )
const (
	Comment  = "--"
	Position = "pos"
	Source   = "src"
)

func IsMarkup(s string) (okay bool) {
	if ch, n := utf8.DecodeRuneInString(s); n > 0 {
		// testing "!IsUpper()" instead of "IsLower()" allows symbols for metadata.
		normal := n < unicode.MaxASCII && unicode.IsUpper(ch)
		okay = !normal
	}
	return
}

// read a user comment from markup, normalizing it as an array of strings
func UserComment(markup map[string]any) (ret []string) {
	switch cmt := markup[Comment].(type) {
	case string:
		ret = []string{cmt}
	case []string:
		ret = cmt
	case []interface{}:
		lines := make([]string, len(cmt))
		for i, el := range cmt {
			if str, ok := el.(string); !ok {
				lines = nil
				break
			} else {
				lines[i] = str
			}
		}
		ret = lines
	}
	return
}
