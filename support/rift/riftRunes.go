package rift

import "unicode"

const (
	InlineArraySeparator = ','
	InlineArrayStop      = '.'
	// interpreted strings are bookended with double quotes
	InterpretedQuotes = '"'
	// values in a sequence are prefixed by a dash ( and whitespace )
	SequenceDash = '-'
	// keywords in a signature are separated by a colon
	SignatureSeparator = ':'
	// valid in words between colons
	// ( as is space and any unicode letter )
	SignatureConnector = '_'
	Space              = ' '
	Newline            = '\n'
	Tab                = '\t' // tab is considered invalid whitespace

)

func isSigWord(r rune) bool {
	return r == SignatureConnector || unicode.IsLetter(r) || r == Space
}
