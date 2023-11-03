package rift

const (
	InlineArraySeparator = ','
	InlineArrayStop      = '.'
	InterpretedString    = '"' // interpreted strings are bookended with double quotes
	Newline              = '\n'
	RawString            = '`'
	SequenceDash         = '-' // values in a sequence are prefixed by a dash ( and whitespace )
	SignatureConnector   = '_' // valid in words between colons
	SignatureSeparator   = ':' // keywords in a signature are separated by a colon
	Space                = ' '
	Tab                  = '\t' // tab is considered invalid whitespace
)
