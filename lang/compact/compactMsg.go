package compact

import "log"

// A decoding of a tapestry command.
// Commands on disk are stored as plain value maps
// containing a command signature and command data pair;
// and any number of markup key-value pairs.
// Signatures start with a capital letter, and contain parameter names separated by colons.
// Markup keys start with a double dash.
type Message struct {
	Key    string         // original specified text: "Sig:label:"
	Lede   string         // the leading part of the signature in lowercase
	Labels []string       // parameter names sans colons; the first can be blank ( anonymous )
	Args   []any          // the same length as labels
	Markup map[string]any // from map keys starting with "--"; stored stripped of the dashes.
}

func (op *Message) AddMarkup(k string, v any) {
	if op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	// filter markup to types that gob can handle
	switch v := v.(type) {
	case bool, string, float64, []int, int:
		op.Markup[k] = v

	case []any:
		if str, ok := SliceStrings(v); !ok {
			// this eats the source pos elements that are encoded into the db
			// they shouldnt even be there really
			log.Printf("unhandled markup %s: %T %v", k, v, v)
			// panic("unhandled markup")
		} else {
			op.Markup[k] = str
		}
	default:
		log.Printf("unhandled markup %s: %T %v", k, v, v)
		panic("unhandled markup")
	}
}
