package js

func Embrace(style [2]rune, cb func(*Builder)) string {
	var out Builder
	return out.Brace(style, cb).String()
}
