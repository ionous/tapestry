package compact

// A decoding of a tapestry command.
// Commands on disk are stored as plain value maps
// containing a command signature and command data pair;
// and any number of markup key-value pairs.
// Signatures start with a capital letter, and contain parameter names separated by colons.
// Markup keys start with a double dash.
type Message struct {
	Key    string         // original specified text: "Sig:label:"
	Name   string         // names are lower_underscore
	Labels []string       // parameter names sans colons
	Args   []any          // the same length as labels
	Markup map[string]any // from map keys starting with "--"; stored stripped of the dashes.
}

func (op *Message) AddMarkup(k string, v any) {
	if op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	op.Markup[k] = v
}
