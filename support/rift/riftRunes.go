package rift

const (
	ArraySeparator     = ','
	ArrayStop          = '.'
	Dash               = '-'  // values in a sequence are prefixed by a dash ( and whitespace )
	Hash               = '#'  // comment marker
	HTab               = '\t' // tab is considered invalid whitespace
	InterpretedString  = '"'  // interpreted strings are bookended with double quotes
	Newline            = '\n'
	RawString          = '`'
	SignatureConnector = '_' // valid in words between colons
	SignatureSeparator = ':' // keywords in a signature are separated by a colon
	Space              = ' '
	Carriage           = '\r' // carriage return is used in comment blocks
)
