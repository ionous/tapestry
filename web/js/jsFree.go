package js

func Embrace(style [2]rune, cb func(*Builder)) string {
	var out Builder
	return out.Brace(style, cb).String()
}

func QuotedStrings(values []string) string {
	var out Builder
	for i, el := range values {
		if i > 0 {
			out.R(Comma)
		}
		out.Q(el)
	}
	return out.String()
}
