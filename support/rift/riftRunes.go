package rift

const (
	ArraySeparator     = ','
	ArrayStop          = '.'
	Dash               = '-'  // values in a sequence are prefixed by a dash ( and whitespace )
	Hash               = '#'  // comment marker
	HTab               = '\t' // tab is considered invalid whitespace
	InterpretedString  = '"'  // interpreted strings are bookended with double quotes
	CollectionMark     = '\r' // in comment blocks, represents the dash or key of a sequence
	Nestline           = '\t' // in comment blocks, indicates nested comments
	Newline            = '\n'
	RawString          = '`'
	Record             = '\f' // form feed is used to separate comment entries
	SignatureConnector = '_'  // valid in words between colons
	SignatureSeparator = ':'  // keywords in a signature are separated by a colon
	Space              = ' '
)
