package assign

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/inflect/en"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
)

// rewrite pattern calls to look like commands
// todo: suss out whether the call pattern name is a literal
// and only then compact it.. that will allow dynamic calls.
func EncodePattern(m jsn.Marshaler, op *CallPattern) (err error) {
	// auto generated command names are underscore separated
	// writeBreak in jsn/cout turns those names into pascal case for the story commands
	// TestEncodePattern checks that common inputs work okay.
	patName := strings.TrimSpace(op.PatternName)
	pb := patternBlock(patName)
	if err = m.MarshalBlock(pb); err == nil {
		for _, arg := range op.Arguments {
			argName := strings.TrimSpace(arg.Name)
			if en.IsCapitalized(argName) {
				argName = en.MixedCaseToSpaces(argName)
			}
			if e := m.MarshalKey(argName, argName); e != nil {
				err = e
				break
			} else if e := rt.Assignment_Marshal(m, &arg.Value); e != nil {
				err = e
				break
			}
		}
		m.EndBlock()
	}
	return
}

// a fake block that writes the pattern name out as the lede
type patternBlock string

func (pb patternBlock) GetLede() string       { return string(pb) }
func (patternBlock) GetType() string          { return "patternBlock" }
func (patternBlock) GetFlow() interface{}     { return nil }
func (patternBlock) SetFlow(interface{}) bool { return false }
