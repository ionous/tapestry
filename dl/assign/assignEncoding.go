package assign

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	inflect "git.sr.ht/~ionous/tapestry/inflect/en"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/rt"
)

func CustomEncoder(enc *encode.Encoder, op jsn.Marshalee) (ret any, err error) {
	if call, ok := op.(*CallPattern); !ok {
		ret, err = literal.CustomEncoder(enc, op)
	} else {
		ret, err = encodePattern(enc, call)
	}
	return
}

// rewrite pattern calls to look like commands
// todo: suss out whether the call pattern name is a literal
// and only then compact it.. that will allow dynamic calls.
func encodePattern(enc *encode.Encoder, op *CallPattern) (ret any, err error) {
	var pb encode.FlowBuilder
	// auto generated command names are underscore separated
	// writeBreak those names into pascal case for the story commands
	// TestEncodePattern checks that common inputs work okay.
	pb.WriteLede(strings.TrimSpace(op.PatternName))
	for _, arg := range op.Arguments {
		argName := strings.TrimSpace(arg.Name)
		if inflect.IsCapitalized(argName) {
			argName = inflect.MixedCaseToSpaces(argName)
		}
		slot := rt.Assignment_Slot{Value: arg.Value}
		if out, e := enc.Encode(&slot); e != nil {
			err = e
			break
		} else {
			pb.WriteArg(argName, out)
		}
	}
	if err == nil {
		ret = pb.FinalizeMap()
	}
	return
}
