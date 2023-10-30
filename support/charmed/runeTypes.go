package charmed

// https://golang.org/ref/spec#decimal_digit
func IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

// https://golang.org/ref/spec#hex_digit
func IsHex(r rune) bool {
	return (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') || IsNumber(r)
}
